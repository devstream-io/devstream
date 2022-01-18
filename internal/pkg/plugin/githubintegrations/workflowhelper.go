package githubintegrations

import (
	"context"
	"fmt"
	"github.com/google/go-github/v40/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

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
