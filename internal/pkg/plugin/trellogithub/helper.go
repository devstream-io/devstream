package trellogithub

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

func buildStatus(tg *TrelloGithub) statemanager.ResourceStatus {
	resStatus := make(statemanager.ResourceStatus)
	resStatus["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", tg.options.Owner, tg.options.Repo)
	return resStatus
}

func (tg *TrelloGithub) buildReadStatus() (statemanager.ResourceStatus, error) {
	listIds := make(statemanager.ResourceStatus)

	path, err := tg.client.GetWorkflowPath()
	if err != nil {
		return nil, err
	}
	listIds["workflowDir"] = path

	return listIds, nil
}
