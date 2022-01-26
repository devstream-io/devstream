package trellogithub

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/google/go-github/v42/github"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/merico-dev/stream/internal/pkg/log"
	gh "github.com/merico-dev/stream/internal/pkg/util/github"
	"github.com/merico-dev/stream/internal/pkg/util/mapz"
	"github.com/merico-dev/stream/internal/pkg/util/slicez"
	"github.com/merico-dev/stream/internal/pkg/util/trello"
)

type TrelloGithub struct {
	ctx     context.Context
	client  *github.Client
	options *Options
}

func NewTrelloGithub(options *map[string]interface{}) (*TrelloGithub, error) {
	ctx := context.Background()

	var opt Options
	err := mapstructure.Decode(*options, &opt)
	if err != nil {
		return nil, err
	}

	if errs := validate(&opt); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return nil, fmt.Errorf("params are illegal")
	}

	client, err := getGitHubClient(ctx)
	if err != nil {
		return nil, err
	}

	return &TrelloGithub{
		ctx:     ctx,
		client:  client,
		options: &opt,
	}, nil
}

func (gi *TrelloGithub) GetApi() *Api {
	return gi.options.Api
}

func (gi *TrelloGithub) AddWorkflow(workflow *Workflow) error {
	sha, err := gi.getFileSHA(workflow.workflowFileName)
	if err != nil {
		return err
	}
	if sha != "" {
		log.Infof("GitHub Actions workflow %s already exists.", workflow.workflowFileName)
		return nil
	}

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(workflow.commitMessage),
		Content: []byte(workflow.workflowContent),
		Branch:  github.String(gi.options.Branch),
	}

	log.Infof("Creating GitHub Actions workflow %s ...", workflow.workflowFileName)
	_, _, err = gi.client.Repositories.CreateFile(
		gi.ctx,
		gi.options.Owner,
		gi.options.Repo,
		generateGitHubWorkflowFileByName(workflow.workflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Infof("Github Actions workflow %s created.", workflow.workflowFileName)
	return nil
}

func (gi *TrelloGithub) DeleteWorkflow(workflow *Workflow) error {
	sha, err := gi.getFileSHA(workflow.workflowFileName)
	if err != nil {
		return err
	}
	if sha == "" {
		log.Successf("Github Actions workflow %s already removed.", workflow.workflowFileName)
		return nil
	}

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(workflow.commitMessage),
		SHA:     github.String(sha),
		Branch:  github.String(gi.options.Branch),
	}

	log.Infof("Deleting GitHub Actions workflow %s ...", workflow.workflowFileName)
	_, _, err = gi.client.Repositories.DeleteFile(
		gi.ctx,
		gi.options.Owner,
		gi.options.Repo,
		generateGitHubWorkflowFileByName(workflow.workflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Successf("GitHub Actions workflow %s removed.", workflow.workflowFileName)
	return nil
}

// VerifyWorkflows get the workflows with names "wf1.yml", "wf2.yml", then:
// If all workflows is ok => return ({"wf1.yml":nil, "wf2.yml:nil}, nil)
// If some error occurred => return (nil, error)
// If wf1.yml is not found => return ({"wf1.yml":error("not found"), "wf2.yml:nil},nil)
func (gi *TrelloGithub) VerifyWorkflows(workflows []*Workflow) (map[string]error, error) {
	wsFiles := make([]string, 0)
	for _, w := range workflows {
		wsFiles = append(wsFiles, w.workflowFileName)
	}

	fmt.Printf("Workflow files: %v", wsFiles)
	filesInRemoteDir, rMap, err := gi.FetchRemoteContent(wsFiles)
	if err != nil {
		return nil, err
	}
	if rMap != nil {
		return rMap, nil
	}

	return gi.CompareFiles(wsFiles, filesInRemoteDir), nil
}

func (gi *TrelloGithub) FetchRemoteContent(wsFiles []string) ([]string, map[string]error, error) {
	var filesInRemoteDir = make([]string, 0)
	_, dirContent, resp, err := gi.client.Repositories.GetContents(
		gi.ctx,
		gi.options.Owner,
		gi.options.Repo,
		".github/workflows",
		&github.RepositoryContentGetOptions{},
	)

	// error reason is not 404
	if err != nil && !strings.Contains(err.Error(), "404") {
		log.Errorf("GetContents failed with error: %s", err)
		return nil, nil, err
	}
	// StatusCode == 404
	if resp.StatusCode == http.StatusNotFound {
		log.Error("GetContents returned with status code 404")
		retMap := mapz.FillMapWithStrAndError(wsFiles, fmt.Errorf("not found"))
		return nil, retMap, nil
	}
	// StatusCode != 200
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("got some error: %s", resp.Status)
	}
	// StatusCode == 200
	log.Info("GetContents return with status code 200")
	for _, f := range dirContent {
		log.Infof("Found remote file: %s", f.GetName())
		filesInRemoteDir = append(filesInRemoteDir, f.GetName())
	}
	return filesInRemoteDir, nil, nil
}

// CompareFiles compare files between local and remote
func (gi *TrelloGithub) CompareFiles(wsFiles, filesInRemoteDir []string) map[string]error {
	lostFiles := slicez.SliceInSliceStr(wsFiles, filesInRemoteDir)
	// all files exist
	if len(lostFiles) == 0 {
		log.Info("All workflows exist.")
		retMap := mapz.FillMapWithStrAndError(wsFiles, nil)
		return retMap
	}
	// some files lost
	retMap := mapz.FillMapWithStrAndError(wsFiles, nil)
	for _, f := range lostFiles {
		log.Warnf("Lost file: %s", f)
		retMap[f] = fmt.Errorf("not found")
	}
	return retMap
}

// getFileSHA will try to collect the SHA hash value of the file, then return it. the return values will be:
// 1. If file exists without error -> string(SHA), nil
// 2. If some errors occurred -> return "", err
// 3. If file not found without error -> return "", nil
func (gi *TrelloGithub) getFileSHA(filename string) (string, error) {
	content, _, resp, err := gi.client.Repositories.GetContents(
		gi.ctx,
		gi.options.Owner,
		gi.options.Repo,
		generateGitHubWorkflowFileByName(filename),
		&github.RepositoryContentGetOptions{},
	)

	// error reason is not 404
	if err != nil && !strings.Contains(err.Error(), "404") {
		return "", err
	}

	// error reason is 404
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	// no error occurred
	if resp.StatusCode == http.StatusOK {
		return *content.SHA, nil
	}
	return "", fmt.Errorf("got some error")
}

// renderTemplate render the github actions template with config.yaml
func (gi *TrelloGithub) renderTemplate(workflow *Workflow) error {
	var jobs Jobs
	err := mapstructure.Decode(gi.options.Jobs, &jobs)
	if err != nil {
		return err
	}
	//if use default {{.}}, it will confict (github actions vars also use them)
	t, err := template.New("trellogithub").Delims("[[", "]]").Parse(workflow.workflowContent)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	err = t.Execute(&buff, jobs)
	if err != nil {
		return err
	}
	workflow.workflowContent = buff.String()
	return nil
}

type TrelloItemId struct {
	boardId     string
	todoListId  string
	doingListId string
	doneListId  string
}

// CreateTrelloItems create board/lists, and set secret by ids
// TODO(daniel-hutao): rename the function name
func (gi *TrelloGithub) CreateTrelloItems() (*TrelloItemId, error) {
	c, _ := trello.NewClient()
	board, err := c.CreateBoard("DevStream_Trello_Board")
	if err != nil {
		return nil, err
	}

	lists, err := board.GetLists()
	if err != nil {
		return nil, err
	}

	var todo string
	var doing string
	var done string

	for _, l := range lists {
		log.Debugf("List name: %s", l.Name)
		if l.Name == "To Do" || l.Name == "待办" {
			todo = l.ID
			continue
		}
		if l.Name == "Doing" || l.Name == "进行中" {
			doing = l.ID
			continue
		}
		if l.Name == "Done" || l.Name == "完成" {
			done = l.ID
			continue
		}
		log.Errorf("Unknow name: %s", l.Name)
		return nil, fmt.Errorf("unknow name: %s", l.Name)
	}

	return &TrelloItemId{
		boardId:     board.ID,
		todoListId:  todo,
		doingListId: doing,
		doneListId:  done,
	}, nil
}

// AddTrelloIdSecret add trello ids to secret
func (gi *TrelloGithub) AddTrelloIdSecret(trelloId *TrelloItemId) error {
	ghOptions := &gh.Option{
		Owner:    gi.options.Owner,
		Repo:     gi.options.Repo,
		NeedAuth: true,
		WorkPath: gh.DefaultWorkPath,
	}
	c, err := gh.NewClient(ghOptions)
	if err != nil {
		return err
	}
	// add key
	if err := c.AddRepoSecret("TRELLO_API_KEY", viper.GetString("trello_api_key")); err != nil {
		return err
	}

	// add token
	if err := c.AddRepoSecret("TRELLO_TOKEN", viper.GetString("trello_token")); err != nil {
		return err
	}

	// add board id
	if err := c.AddRepoSecret("TRELLO_BOARD_ID", trelloId.boardId); err != nil {
		return err
	}

	// add todolist id
	if err := c.AddRepoSecret("TRELLO_TODO_LIST_ID", trelloId.todoListId); err != nil {
		return err
	}

	// add doinglist id
	if err := c.AddRepoSecret("TRELLO_DOING_LIST_ID", trelloId.doingListId); err != nil {
		return err
	}

	// add donelist id
	if err := c.AddRepoSecret("TRELLO_DONE_LIST_ID", trelloId.doneListId); err != nil {
		return err
	}

	// add member map
	if err := c.AddRepoSecret("TRELLO_MEMBER_MAP", "[]"); err != nil {
		return err
	}

	return nil
}
