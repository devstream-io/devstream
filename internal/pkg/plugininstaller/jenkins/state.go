package jenkins

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

func GetStatus(options plugininstaller.RawOptions) (statemanager.ResourceState, error) {
	opts, err := newJobOptions(options)
	if err != nil {
		return nil, err
	}

	client, err := jenkins.NewClient(opts.Jenkins.URL, opts.BasicAuth)
	if err != nil {
		return nil, err
	}
	res := make(statemanager.ResourceState)
	job, err := client.GetFolderJob(
		opts.Pipeline.getJobName(), opts.Pipeline.getJobFolder(),
	)
	if err != nil {
		return nil, err
	}
	rawJob := job.Raw
	res["JobCreated"] = true
	res["Job"] = map[string]interface{}{
		"Created": true,
		"Class":   rawJob.Class,
		"URL":     rawJob.URL,
	}
	return res, nil
}
