package trellogithub

import (
	"github.com/merico-dev/stream/internal/pkg/log"
)

// Reinstall remove and set up trello-github-integ workflows.
func Reinstall(options *map[string]interface{}) (bool, error) {
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

		if err := gis.renderTemplate(w); err != nil {
			return false, err
		}
		err = gis.AddWorkflow(w)
		if err != nil {
			return false, err
		}
	}
	log.Success("Adding workflow file succeeded.")
	trelloIds, err := gis.CreateTrelloItems()
	if err != nil {
		return false, err
	}
	log.Success("Creating trello board succeeded.")
	if err := gis.AddTrelloIdSecret(trelloIds); err != nil {
		return false, err
	}
	log.Success("Adding secret keys for trello succeeded.")
	return true, nil
}
