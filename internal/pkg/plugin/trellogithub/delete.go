package trellogithub

import (
	"github.com/merico-dev/stream/internal/pkg/log"
)

// Delete remove trello-github-integ workflows.
func Delete(options *map[string]interface{}) (bool, error) {
	gis, err := NewTrelloGithub(options)
	if err != nil {
		return false, err
	}

	api := gis.GetApi()
	log.Infof("api is %s", api.Name)
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(api.Name)

	for _, w := range ws {
		err := gis.DeleteWorkflow(w)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
