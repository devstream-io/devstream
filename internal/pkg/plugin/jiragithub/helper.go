package jiragithub

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/log"
)

func parseAndValidateOptions(options map[string]interface{}) (*Options, error) {
	var opt Options
	err := mapstructure.Decode(options, &opt)
	if err != nil {
		return nil, err
	}

	if errs := validate(&opt); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return nil, fmt.Errorf("incorrect params")
	}

	return &opt, nil
}

func setRepoSecrets(gitHubClient *github.Client) error {

	// JIRA_API_TOKEN, how to get it: "https://help.siteimprove.com/support/solutions/articles/80000448174-how-to-create-an-api-token-from-your-atlassian-account"
	if err := gitHubClient.AddRepoSecret("JIRA_API_TOKEN", viper.GetString("jira_api_token")); err != nil {
		return err
	}

	if err := gitHubClient.AddRepoSecret("GH_TOKEN", viper.GetString("github_token")); err != nil {
		return err
	}

	return nil
}
