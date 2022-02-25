package trellogithub

import (
	"github.com/merico-dev/stream/pkg/util/log"
)

// Update remove and set up trello-github-integ workflows.
func Update(options *map[string]interface{}) (map[string]interface{}, error) {
	gis, err := NewTrelloGithub(options)
	if err != nil {
		return nil, err
	}

	api := gis.GetApi()
	log.Infof("API is %s.", api.Name)
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(api.Name)

	for _, w := range ws {
		err := gis.DeleteWorkflow(w)
		if err != nil {
			return nil, err
		}

		if err := gis.renderTemplate(w); err != nil {
			return nil, err
		}
		err = gis.AddWorkflow(w)
		if err != nil {
			return nil, err
		}
	}
	log.Success("Adding workflow file succeeded.")
	trelloIds, err := gis.CreateTrelloItems()
	if err != nil {
		return nil, err
	}
	log.Success("Creating trello board succeeded.")
	if err := gis.AddTrelloIdSecret(trelloIds); err != nil {
		return nil, err
	}

	log.Success("Adding secret keys for trello succeeded.")

	return buildState(gis, trelloIds), nil
}
