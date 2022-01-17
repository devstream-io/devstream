package githubintegrations

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/spf13/viper"

	"github.com/google/go-github/v40/github"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"

	"github.com/merico-dev/stream/internal/pkg/util/mapz"
	"github.com/merico-dev/stream/internal/pkg/util/slicez"
)

type GithubIntegrations struct {
	ctx     context.Context
	client  *github.Client
	options *Options
}

func NewGithubIntegrations(options *map[string]interface{}) (*GithubIntegrations, error) {
	ctx := context.Background()

	var opt Options
	err := mapstructure.Decode(*options, &opt)
	if err != nil {
		return nil, err
	}

	if !verifyOptions(&opt) {
		return nil, fmt.Errorf("options is illegal")
	}

	client, err := getGitHubClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GithubIntegrations{
		ctx:     ctx,
		client:  client,
		options: &opt,
	}, nil
}

func (gi *GithubIntegrations) GetApi() *Api {
	return gi.options.Api
}

func (gi *GithubIntegrations) AddWorkflow(workflow *Workflow) error {
	sha, err := gi.getFileSHA(workflow.workflowFileName)
	if err != nil {
		return err
	}
	if sha != "" {
		log.Printf("GitHub Actions workflow %s already exists.", workflow.workflowFileName)
		return nil
	}

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(workflow.commitMessage),
		Content: []byte(workflow.workflowContent),
		Branch:  github.String(gi.options.Branch),
	}

	log.Printf("Creating GitHub Actions workflow %s ...", workflow.workflowFileName)
	_, _, err = gi.client.Repositories.CreateFile(
		gi.ctx,
		gi.options.Owner,
		gi.options.Repo,
		generateGitHubWorkflowFileByName(workflow.workflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Printf("Github Actions workflow %s created.", workflow.workflowFileName)
	return nil
}

func (gi *GithubIntegrations) DeleteWorkflow(workflow *Workflow) error {
	sha, err := gi.getFileSHA(workflow.workflowFileName)
	if err != nil {
		return err
	}
	if sha == "" {
		log.Printf("Github Actions workflow %s already removed.", workflow.workflowFileName)
		return nil
	}

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(workflow.commitMessage),
		SHA:     github.String(sha),
		Branch:  github.String(gi.options.Branch),
	}

	log.Printf("Deleting GitHub Actions workflow %s ...", workflow.workflowFileName)
	_, _, err = gi.client.Repositories.DeleteFile(
		gi.ctx,
		gi.options.Owner,
		gi.options.Repo,
		generateGitHubWorkflowFileByName(workflow.workflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Printf("GitHub Actions workflow %s removed.", workflow.workflowFileName)
	return nil
}

// VerifyWorkflows get the workflows with names "wf1.yml", "wf2.yml", then:
// If all workflows is ok => return ({"wf1.yml":nil, "wf2.yml:nil}, nil)
// If some error occurred => return (nil, error)
// If wf1.yml is not found => return ({"wf1.yml":error("not found"), "wf2.yml:nil},nil)
func (gi *GithubIntegrations) VerifyWorkflows(workflows []*Workflow) (map[string]error, error) {
	wsFiles := make([]string, 0)
	for _, w := range workflows {
		wsFiles = append(wsFiles, w.workflowFileName)
	}
	fmt.Printf("Workflow files: %v", wsFiles)

	_, dirContent, resp, err := gi.client.Repositories.GetContents(
		gi.ctx,
		gi.options.Owner,
		gi.options.Repo,
		".github/workflows",
		&github.RepositoryContentGetOptions{},
	)

	// error reason is not 404
	if err != nil && !strings.Contains(err.Error(), "404") {
		log.Printf("GetContents failed with error: %s", err)
		return nil, err
	}
	// StatusCode == 404
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("GetContents return with status code 404")
		retMap := mapz.FillMapWithStrAndError(wsFiles, fmt.Errorf("not found"))
		return retMap, nil
	}
	// StatusCode != 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got some error is not expected: %s", resp.Status)
	}
	// StatusCode == 200
	log.Printf("GetContents return with status code 200")
	var filesInRemoteDir = make([]string, 0)
	for _, f := range dirContent {
		log.Printf("Found remote file: %s", f.GetName())
		filesInRemoteDir = append(filesInRemoteDir, f.GetName())
	}

	lostFiles := slicez.SliceInSliceStr(wsFiles, filesInRemoteDir)
	// all files exist
	if len(lostFiles) == 0 {
		log.Println("All files exist")
		retMap := mapz.FillMapWithStrAndError(wsFiles, nil)
		return retMap, nil
	}
	// some files lost
	log.Println("Some files lost")
	retMap := mapz.FillMapWithStrAndError(wsFiles, nil)
	for _, f := range lostFiles {
		log.Printf("Lost file: %s", f)
		retMap[f] = fmt.Errorf("not found")
	}
	return retMap, nil
}

// getFileSHA will try to collect the SHA hash value of the file, then return it. the return values will be:
// 1. If file exists without error -> string(SHA), nil
// 2. If some errors occurred -> return "", err
// 3. If file not found without error -> return "", nil
func (gi *GithubIntegrations) getFileSHA(filename string) (string, error) {
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
	return "", fmt.Errorf("got some error is not expected")
}

// renderTemplate render the github actions template with config.yaml
func (gi *GithubIntegrations) renderTemplate(workflow *Workflow) error {
	var jobs Jobs
	err := mapstructure.Decode(gi.options.Jobs, &jobs)
	if err != nil {
		return err
	}
	//if use default {{.}}, it will confict (github actions vars also use them)
	t, err := template.New("githubactions").Delims("[[", "]]").Parse(workflow.workflowContent)
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

func generateGitHubWorkflowFileByName(f string) string {
	return fmt.Sprintf(".github/workflows/%s", f)
}

func getGitHubClient(ctx context.Context) (*github.Client, error) {
	token := viper.GetString("github_token")
	if token == "" {
		return nil, fmt.Errorf("failed to initialize GitHub token. More info - https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

func verifyOptions(opt *Options) bool {
	return opt.Owner != "" &&
		opt.Repo != "" &&
		opt.Branch != "" &&
		opt.Api != nil
}
