package trellogithub

import (
	"bytes"
	"fmt"

	"text/template"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/github"
	"github.com/merico-dev/stream/pkg/util/trello"
)

// renderTemplate render the github actions template with config.yaml
func (tg *TrelloGithub) renderTemplate(workflow *github.Workflow) error {
	var jobs Jobs
	err := mapstructure.Decode(tg.options.Jobs, &jobs)
	if err != nil {
		return err
	}
	//if use default {{.}}, it will confict (github actions vars also use them)
	t, err := template.New("trellogithub").Delims("[[", "]]").Parse(workflow.WorkflowContent)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	err = t.Execute(&buff, jobs)
	if err != nil {
		return err
	}
	workflow.WorkflowContent = buff.String()
	return nil
}

func buildState(tg *TrelloGithub, ti *TrelloItemId) map[string]interface{} {
	res := make(map[string]interface{})
	res["workflowDir"] = fmt.Sprintf("/repos/%s/%s/contents/.github/workflows", tg.options.Owner, tg.options.Repo)
	res["boardId"] = ti.boardId
	res["todoListId"] = ti.todoListId
	res["doingListId"] = ti.doingListId
	res["doneListId"] = ti.doneListId
	return res
}

func (tg *TrelloGithub) buildReadState(api *Api) (map[string]interface{}, error) {
	c, err := trello.NewClient()
	if err != nil {
		return nil, err
	}
	listIds, err := c.GetBoardIdAndListId(tg.options.Owner, tg.options.Repo, api.KanbanBoardName)
	if err != nil {
		return nil, err
	}

	path, err := tg.GetWorkflowPath()
	if err != nil {
		return nil, err
	}
	listIds["workflowDir"] = path
	return listIds, nil
}
