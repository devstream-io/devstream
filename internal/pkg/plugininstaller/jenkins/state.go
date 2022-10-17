package jenkins

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

func GetStatus(options plugininstaller.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}

	client, err := jenkins.NewClient(opts.Jenkins.URL, opts.BasicAuth)
	if err != nil {
		return nil, err
	}

	res := make(statemanager.ResourceStatus)
	jobRes, err := getJobState(client, opts.Pipeline.getJobName(), opts.Pipeline.getJobFolder())
	if err != nil {
		return nil, err
	}
	res["JobCreated"] = true
	res["Job"] = jobRes
	return res, nil
}

func getJobState(jenkinsClient jenkins.JenkinsAPI, jobName, jobFolder string) (map[string]interface{}, error) {
	job, err := jenkinsClient.GetFolderJob(jobName, jobFolder)
	if err != nil {
		return nil, err
	}
	rawJob := job.Raw
	return map[string]interface{}{
		"Created": true,
		"Class":   rawJob.Class,
		"URL":     rawJob.URL,
	}, nil
}
