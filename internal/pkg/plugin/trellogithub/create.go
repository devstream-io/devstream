package trellogithub

import (
	"github.com/merico-dev/stream/pkg/util/log"
)

// Create sets up trello-github-integ workflows.
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	tg, err := NewTrelloGithub(options)
	if err != nil {
		return nil, err
	}

	api := tg.GetApi()
	log.Infof("API is: %s.", api.Name)
	ws := defaultWorkflows.GetWorkflowByNameVersionTypeString(api.Name)

	for _, w := range ws {
		if err := tg.renderTemplate(w); err != nil {
			return nil, err
		}
		if err := tg.client.AddWorkflow(w, tg.options.Branch); err != nil {
			return nil, err
		}
	}
	log.Success("Adding workflow file succeeded.")

	trelloIds := &TrelloItemId{
		boardId:     tg.GetApi().BoardId,
		todoListId:  tg.GetApi().todoListId,
		doingListId: tg.GetApi().doingListId,
		doneListId:  tg.GetApi().doneListId,
	}

	if err := tg.AddTrelloIdSecret(trelloIds); err != nil {
		return nil, err
	}

	log.Success("Adding secret keys for trello succeeded.")

	return buildState(tg, trelloIds), nil
}
