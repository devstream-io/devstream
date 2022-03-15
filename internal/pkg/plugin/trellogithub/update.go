package trellogithub

import (
	"github.com/merico-dev/stream/pkg/util/log"
)

// Update remove and set up trello-github-integ workflows.
func Update(options map[string]interface{}) (map[string]interface{}, error) {
	tg, err := NewTrelloGithub(options)
	if err != nil {
		return nil, err
	}

	err = tg.client.DeleteWorkflow(trelloWorkflow, tg.options.Branch)
	if err != nil {
		return nil, err
	}

	err = tg.client.AddWorkflow(trelloWorkflow, tg.options.Branch)
	if err != nil {
		return nil, err
	}

	log.Success("Adding workflow file succeeded.")

	trelloIds := &TrelloItemId{
		boardId:     tg.options.BoardId,
		todoListId:  tg.options.todoListId,
		doingListId: tg.options.doingListId,
		doneListId:  tg.options.doneListId,
	}

	if err := tg.AddTrelloIdSecret(trelloIds); err != nil {
		return nil, err
	}
	log.Success("Adding secret keys for trello succeeded.")

	return buildState(tg, trelloIds), nil
}
