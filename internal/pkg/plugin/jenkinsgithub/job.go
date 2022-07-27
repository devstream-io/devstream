package jenkinsgithub

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

//go:embed job-pr.xml
var jobPrTemplate string

func createJob(client *jenkins.Jenkins, jobName, jobTemplate string, opts *Options) error {
	jobOpts := &JobOptions{
		JobName:              jobName,
		PrPipelineScriptPath: opts.J.PipelineScriptPath,
		GitHubRepoURL:        opts.GitHubRepoURL,
		AdminList:            opts.AdminList,
		CredentialsID:        jenkinsCredentialID,
	}

	jobContent, err := renderJobXml(jobTemplate, jobOpts)
	if err != nil {
		return fmt.Errorf("failed to render job xml: %s", err)
	}

	if _, err := client.CreateJob(context.Background(), jobContent, jobName); err != nil {
		return fmt.Errorf("failed to create job: %s", err)
	}

	return nil
}

func renderJobXml(jobTemplate string, opts *JobOptions) (string, error) {
	tpl := template.New(githubIntegName)
	tpl, err := tpl.Parse(jobTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, opts)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
