package githubactions

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	sha, err := ga.getFileSHA(workflow.workflowFileName)
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
		Branch:  github.String(ga.options.Branch),
	}

	log.Printf("Creating GitHub Actions workflow %s ...", workflow.workflowFileName)
	_, _, err = ga.client.Repositories.CreateFile(
		ga.ctx,
		ga.options.Owner,
		ga.options.Repo,
		generateGitHubWorkflowFileByName(workflow.workflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Printf("Github Actions workflow %s created.", workflow.workflowFileName)
	return nil
}

func (ga *GithubActions) DeleteWorkflow(workflow *Workflow) error {
	sha, err := ga.getFileSHA(workflow.workflowFileName)
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
		Branch:  github.String(ga.options.Branch),
	}

	log.Printf("Deleting GitHub Actions workflow %s ...", workflow.workflowFileName)
	_, _, err = ga.client.Repositories.DeleteFile(
		ga.ctx,
		ga.options.Owner,
		ga.options.Repo,
		generateGitHubWorkflowFileByName(workflow.workflowFileName),
		opts)

	if err != nil {
		return err
	}
	log.Printf("GitHub Actions workflow %s removed.", workflow.workflowFileName)
	return nil
}

// getFileSHA will try to collect the SHA hash value of the file, then return it. the return values will be:
// 1. If file exists without error -> string(SHA), nil
// 2. If some errors occurred -> return "", err
// 3. If file not found without error -> return "", nil
func (ga *GithubActions) getFileSHA(filename string) (string, error) {
	content, _, resp, err := ga.client.Repositories.GetContents(
		ga.ctx,
		ga.options.Owner,
		ga.options.Repo,
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
