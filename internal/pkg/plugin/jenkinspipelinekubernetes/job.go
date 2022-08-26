package jenkinspipelinekubernetes

import (
	"context"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func CreateJob(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// get the jenkins client and test the connection
	client, err := NewJenkinsFromOptions(opts)
	if err != nil {
		return err
	}

	// create credential if not exists
	if _, err := client.GetCredentialsUsername(jenkinsCredentialID); err != nil {
		log.Infof("credential %s not found, creating...", jenkinsCredentialID)
		if err := client.CreateCredentialsUsername(jenkinsCredentialUsername, githubToken, jenkinsCredentialID, jenkinsCredentialDesc); err != nil {
			return err
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
			return err
		}
		if _, err := client.CreateJob(ctx, jobXmlContent, opts.JobName); err != nil {
			return fmt.Errorf("failed to create job: %s", err)
		}
	}

	return nil
}

func DeleteJob(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// get the jenkins client and test the connection
	client, err := NewJenkinsFromOptions(opts)
	if err != nil {
		return err
	}

	// delete the credentials created by devstream if exists
	if _, err := client.GetCredentialsUsername(jenkinsCredentialID); err == nil {
		if err := client.DeleteCredentialsUsername(jenkinsCredentialID); err != nil {
			return err
		}
	}

	// delete the job created by devstream if exists
	if _, err = client.GetJob(context.Background(), opts.JobName); err == nil {
		if _, err := client.DeleteJob(context.Background(), opts.JobName); err != nil {
			return err
		}
	}

	return nil
}

func GetState(options plugininstaller.RawOptions) (statemanager.ResourceState, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	// get the jenkins client and test the connection
	client, err := NewJenkinsFromOptions(opts)
	if err != nil {
		return nil, err
	}

	res := &resource{}

	if _, err = client.GetCredentialsUsername(jenkinsCredentialID); err == nil {
		res.CredentialsCreated = true
	}

	if _, err = client.GetJob(context.Background(), opts.JobName); err == nil {
		res.JobCreated = true
	}

	return res.toMap(), nil
}
