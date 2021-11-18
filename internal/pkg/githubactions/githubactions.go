package githubactions

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v40/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func generateGitHubWorkflowFileByName(f string) string {
	return fmt.Sprintf("%s/%s", ".github/workflows", f)
}

func getGitHubToken() string {
	viper.BindEnv("github_token")
	token, ok := viper.Get("github_token").(string)

	if !ok {
		log.Fatalf("ENV var GITHUB_TOKEN is needed")
	}
	return token
}

func getGitHubClient(ctx *context.Context) *github.Client {
	token := getGitHubToken()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(*ctx, ts)
	return github.NewClient(tc)
}

func fileExists(params *Param) bool {
	_, _, resp, err := params.client.Repositories.GetContents(
		*params.ctx,
		params.options.Owner,
		params.options.Repo,
		generateGitHubWorkflowFileByName(params.pipeline.workflowFileName),
		&github.RepositoryContentGetOptions{})

	if err != nil {
		if resp.StatusCode == 401 {
			log.Fatal("invalid GitHub credentials in GITHUB_TOKEN ENV var")
		} else if resp.StatusCode != 404 {
			log.Fatal(err)
		}
	}

	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

func createFile(params *Param) {
	if fileExists(params) {
		log.Printf("github actions workflow %s already exists\n", params.pipeline.workflowFileName)
		return
	}

	// Note: the file needs to be absent from the repository as you are not
	// specifying a SHA reference here.
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(params.pipeline.commitMessage),
		Content: []byte(params.pipeline.workflowContent),
		Branch:  github.String("master"),
	}

	log.Printf("creating github actions workflow %s...\n", params.pipeline.workflowFileName)
	_, _, err := params.client.Repositories.CreateFile(
		*params.ctx,
		params.options.Owner,
		params.options.Repo,
		generateGitHubWorkflowFileByName(params.pipeline.workflowFileName),
		opts)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Printf("github actions workflow %s created\n", params.pipeline.workflowFileName)
}
