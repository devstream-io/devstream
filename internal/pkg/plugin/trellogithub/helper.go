package trellogithub

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/google/go-github/v42/github"
	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/trello"
)

// getFileSHA will try to collect the SHA hash value of the file, then return it. the return values will be:
// 1. If file exists without error -> string(SHA), nil
// 2. If some errors occurred -> return "", err
// 3. If file not found without error -> return "", nil
func (gi *TrelloGithub) getFileSHA(filename string) (string, error) {
	content, _, resp, err := gi.client.Repositories.GetContents(
		gi.ctx,
		gi.options.Owner,
		gi.options.Repo,
		generateGitHubWorkflowFileByName(filename),
		&github.RepositoryContentGetOptions{},
	)

	// error reason is not 404
	if err != nil && !strings.Contains(err.Error(), "404") {
		return "", err
	}

	// error reason is 404
	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	// no error occurred
	if resp.StatusCode == http.StatusOK {
		return *content.SHA, nil
	}
	return "", fmt.Errorf("got some error")
}

// renderTemplate render the github actions template with config.yaml
func (gi *TrelloGithub) renderTemplate(workflow *Workflow) error {
	var jobs Jobs
	err := mapstructure.Decode(gi.options.Jobs, &jobs)
	if err != nil {
		return err
	}
	//if use default {{.}}, it will confict (github actions vars also use them)
	t, err := template.New("trellogithub").Delims("[[", "]]").Parse(workflow.workflowContent)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	err = t.Execute(&buff, jobs)
	if err != nil {
		return err
	}
	workflow.workflowContent = buff.String()
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

func (gi *TrelloGithub) buildReadState(api *Api) (map[string]interface{}, error) {
	c, err := trello.NewClient()
	if err != nil {
		return nil, err
	}
	listIds, err := c.GetBoardIdAndListId(gi.options.Owner, gi.options.Repo, api.KanbanBoardName)
	if err != nil {
		return nil, err
	}

	path, err := gi.GetWorkflowPath()
	if err != nil {
		return nil, err
	}
	listIds["workflowDir"] = path
	return listIds, nil
}
