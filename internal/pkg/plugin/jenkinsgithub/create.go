package jenkinsgithub

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validateAndHandleOptions(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	// get the jenkins client and test the connection
	client, err := NewJenkinsFromOptions(&opts)
	if err != nil {
		return nil, err
	}

	// create credential if not exists
	if _, err := client.GetCredentialsUsername(jenkinsCredentialID); err != nil {
		log.Infof("credential %s not found, creating...", jenkinsCredentialID)
		if err := client.CreateCredentialsUsername(jenkinsCredentialUsername, opts.GitHubToken, jenkinsCredentialID, jenkinsCredentialDesc); err != nil {
			return nil, err
		}
	}

	// install plugin if not exists
	if err := installPluginsIfNotExists(client); err != nil {
		return nil, err
	}

	// set config about GitHub integration
	if err := applyGitHubIntegConfig(&opts); err != nil {
		return nil, fmt.Errorf("failed to apply github integ config: %s", err)
	}

	// create a job triggered by GitHub pull request and auto create GitHub webhooks
	if err := createJob(client, opts.J.JobName, jobPrTemplate, &opts); err != nil {
		return nil, fmt.Errorf("failed to create job: %s", err)
	}

	// TODO(aFlyBird0): create a job triggered when commits pushed to main(master) branch

	return Read(options)
}
