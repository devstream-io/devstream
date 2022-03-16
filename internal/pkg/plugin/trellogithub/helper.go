package trellogithub

import (
	"fmt"
)

func buildState(tg *TrelloGithub) map[string]interface{} {
	res := make(map[string]interface{})
	res["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", tg.options.Owner, tg.options.Repo)
	return res
}

func (tg *TrelloGithub) buildReadState() (map[string]interface{}, error) {
	listIds := make(map[string]interface{})

	path, err := tg.client.GetWorkflowPath()
	if err != nil {
		return nil, err
	}
	listIds["workflowDir"] = path

	return listIds, nil
}
