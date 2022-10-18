package github

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func GetStaticWorkFlowStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := NewGithubActionOptions(options)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	if opts.Owner != "" {
		res["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", opts.Owner, opts.Repo)
	} else {
		res["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", opts.Org, opts.Repo)
	}
	return res, nil
}

func GetActionStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := NewGithubActionOptions(options)
	if err != nil {
		return nil, err
	}

	log.Debugf("Language is: %s.", opts.GetLanguage())
	ghClient, err := opts.GetGithubClient()
	if err != nil {
		return nil, err
	}
	path, err := ghClient.GetWorkflowPath()
	if err != nil {
		return nil, err
	}

	if path == "" {
		// file not found
		log.Debug("Github action file not found")
		return nil, nil
	}

	resStatus := make(statemanager.ResourceStatus)
	resStatus["workflowDir"] = path
	return resStatus, nil
}
