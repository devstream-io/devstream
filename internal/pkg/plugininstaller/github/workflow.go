package github

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/github"
)

func BuildWorkFlowsWrapper(workflows []*github.Workflow) plugininstaller.MutableOperation {
	return func(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
		opts, err := NewGithubActionOptions(options)
		if err != nil {
			return options, err
		}
		for _, w := range workflows {
			content, err := opts.RenderWorkFlow(w.WorkflowContent)
			if err != nil {
				return options, err
			}
			w.WorkflowContent = content
			opts.Workflows = append(opts.Workflows, w)
		}
		return opts.Encode()
	}
}
