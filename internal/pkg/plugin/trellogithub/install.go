package trellogithub

import (
	"github.com/merico-dev/stream/internal/pkg/log"
)

// Install sets up trello-github-integ workflows.
func Install(options *map[string]interface{}) (bool, error) {
	gis, err := NewTrelloGithub(options)
	if err != nil {
		return false, err
	}

	api := gis.GetApi()
	log.Infof("api is: %s.", api.Name)
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(api.Name)

	for _, w := range ws {
		if err := gis.renderTemplate(w); err != nil {
			return false, err
		}
		if err := gis.AddWorkflow(w); err != nil {
			return false, err
		}
	}
	log.Success("Succeeded in adding workflow file success.")
	trelloIds, err := gis.CreateTrelloItems()
	if err != nil {
		return false, err
	}
	log.Success("Succeeded in creating trello board.")
	if err := gis.AddTrelloIdSecret(trelloIds); err != nil {
		return false, err
	}
	log.Success("Succeeded in adding secret keys for trello.")
	return true, nil
}
