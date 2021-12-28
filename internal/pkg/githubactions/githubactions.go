package githubactions

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v40/github"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type GithubActions struct {
	ctx     context.Context
	client  *github.Client
	options *Options
}

func NewGithubActions(options *map[string]interface{}) (*GithubActions, error) {
	ctx := context.Background()

	var opt Options
	err := mapstructure.Decode(*options, &opt)
	if err != nil {
		return nil, err
	}

	client, err := getGitHubClient(ctx)
	if err != nil {
		return nil, err
	}

	return &GithubActions{
		ctx:     ctx,
		client:  client,
		options: &opt,
	}, nil
}

func (ga *GithubActions) GetLanguage() *Language {
	return ga.options.Language
}

func (ga *GithubActions) AddWorkflow(workflow *Workflow) error {
	exists, err := ga.fileExists(workflow.workflowFileName)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("github actions Workflow %s already exists\n", workflow.workflowFileName)
		return nil
	}

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(workflow.commitMessage),
		Content: []byte(workflow.workflowContent),
		Branch:  github.String(ga.options.Branch),
	}

	log.Printf("creating github actions Workflow %s...\n", workflow.workflowFileName)
	_, _, err = ga.client.Repositories.CreateFile(
		ga.ctx,
		ga.options.Owner,
		ga.options.Repo,
		generateGitHubWorkflowFileByName(workflow.workflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Printf("github actions Workflow %s created\n", workflow.workflowFileName)
	return nil
}

func (ga *GithubActions) DeleteWorkflow(workflow *Workflow) error {
	exists, err := ga.fileExists(workflow.workflowFileName)
	if err != nil {
		return err
	}
	if !exists {
		log.Printf("github actions Workflow %s already removed\n", workflow.workflowFileName)
		return nil
	}

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(workflow.commitMessage),
		Content: []byte(workflow.workflowContent),
		Branch:  github.String(ga.options.Branch),
	}

	log.Printf("deleting github actions Workflow %s...\n", workflow.workflowFileName)
	_, _, err = ga.client.Repositories.DeleteFile(
		ga.ctx,
		ga.options.Owner,
		ga.options.Repo,
		generateGitHubWorkflowFileByName(workflow.workflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Printf("github actions Workflow %s removed\n", workflow.workflowFileName)
	return nil
}

func (ga *GithubActions) fileExists(filename string) (bool, error) {
	_, _, resp, err := ga.client.Repositories.GetContents(
		ga.ctx,
		ga.options.Owner,
		ga.options.Repo,
		generateGitHubWorkflowFileByName(filename),
		&github.RepositoryContentGetOptions{},
	)

	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, fmt.Errorf("got some error is not expected")
}

func generateGitHubWorkflowFileByName(f string) string {
	return fmt.Sprintf(".github/workflows/%s", f)
}

func getGitHubToken() string {
	err := viper.BindEnv("github_token")
	if err != nil {
		log.Printf("bind ENV var GITHUB_TOKEN failed: %s", err)
		return ""
	}

	return viper.GetString("github_token")
}

func getGitHubClient(ctx context.Context) (*github.Client, error) {
	token := getGitHubToken()
	if token == "" {
		return nil, fmt.Errorf("failed to initialize GitHub token. More info - https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}
