package trellogithub

import (
	"github.com/merico-dev/stream/pkg/util/log"
)

// Delete remove trello-github-integ workflows.
func Delete(options map[string]interface{}) (bool, error) {
	tg, err := NewTrelloGithub(options)
	if err != nil {
		return false, err
	}

	api := tg.GetApi()
	log.Infof("API is %s.", api.Name)
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(api.Name)

	for _, w := range ws {
		err := tg.client.DeleteWorkflow(w, tg.options.Branch)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
