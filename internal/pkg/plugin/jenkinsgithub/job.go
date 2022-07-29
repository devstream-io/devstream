package jenkinsgithub

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/template"
)

//go:embed tpl/job-pr.tpl.xml
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

// renderJobXml renders the job xml content from the job template and the job options.
// pr job && main job
func renderJobXml(jobTemplate string, opts *JobOptions) (string, error) {
	return template.Render(githubIntegName, jobTemplate, opts)
}
