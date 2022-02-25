package trellogithub

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/spf13/viper"

	gh "github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
	"github.com/merico-dev/stream/pkg/util/mapz"
)

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
		log.Errorf("GetContents failed with error: %s.", err)
		return nil, nil, err
	}
	// StatusCode == 404
	if resp.StatusCode == http.StatusNotFound {
		log.Error("GetContents returned with status code 404.")
		retMap := mapz.FillMapWithStrAndError(wsFiles, fmt.Errorf("not found"))
		return nil, retMap, nil
	}
	// StatusCode != 200
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("got some error: %s", resp.Status)
	}
	// StatusCode == 200
	log.Info("GetContents return with status code 200.")
	for _, f := range dirContent {
		log.Infof("Found remote file: %s.", f.GetName())
		filesInRemoteDir = append(filesInRemoteDir, f.GetName())
	}
	return filesInRemoteDir, nil, nil
}

// AddTrelloIdSecret add trello ids to secret
func (gi *TrelloGithub) AddTrelloIdSecret(trelloId *TrelloItemId) error {
	ghOptions := &gh.Option{
		Owner:    gi.options.Owner,
		Repo:     gi.options.Repo,
		NeedAuth: true,
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

func (gi *TrelloGithub) GetWorkflowPath() (string, error) {
	_, _, resp, err := gi.client.Repositories.GetContents(
		gi.ctx,
		gi.options.Owner,
		gi.options.Repo,
		".github/workflows",
		&github.RepositoryContentGetOptions{},
	)

	// error reason is not 404
	if err != nil && !strings.Contains(err.Error(), "404") {
		return "", err
	}

	// error reason is 404
	if resp.StatusCode == http.StatusNotFound {
		return "", err
	}

	return resp.Request.URL.Path, nil
}
