package jenkinspipelinekubernetes

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

const (
	jenkinsCredentialID       = "credential-jenkins-pipeline-kubernetes-by-devstream"
	jenkinsCredentialDesc     = "Jenkins Pipeline secret, created by devstream/jenkins-pipeline-kubernetes"
	jenkinsCredentialUsername = "foo-useless-username"
)

var githubToken string

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
		if err := client.CreateCredentialsUsername(jenkinsCredentialUsername, githubToken, jenkinsCredentialID, jenkinsCredentialDesc); err != nil {
			return nil, err
		}
	}

	// create job if not exists
	ctx := context.Background()
	if _, err := client.GetJob(ctx, opts.JobName); err != nil {
		log.Infof("job %s not found, creating...", opts.JobName)
		jobXmlOpts := &JobXmlOptions{
			GitHubRepoURL:      opts.JenkinsfileScmURL,
			CredentialsID:      jenkinsCredentialID,
			PipelineScriptPath: opts.JenkinsfilePath,
		}
		jobXmlContent, err := renderJobXml(jobTemplate, jobXmlOpts)
		if err != nil {
			return nil, err
		}
		if _, err := client.CreateJob(ctx, jobXmlContent, opts.JobName); err != nil {
			return nil, fmt.Errorf("failed to create job: %s", err)
		}
	}

	res := &resource{
		CredentialsCreated: true,
		JobCreated:         true,
	}

	return res.toMap(), nil
}

type JobXmlOptions struct {
	GitHubRepoURL      string
	CredentialsID      string
	PipelineScriptPath string
}

func renderJobXml(jobTemplate string, opts *JobXmlOptions) (string, error) {
	return template.Render("jenkins-job", jobTemplate, opts)
}
