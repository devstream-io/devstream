package trellogithub

import (
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create sets up trello-github-integ workflows.
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	tg, err := NewTrelloGithub(options)
	if err != nil {
		return nil, err
	}

	if err := tg.client.AddWorkflow(trelloWorkflow, tg.options.Branch); err != nil {
		return nil, err
	}

	log.Success("Adding workflow file succeeded.")

	trelloItemId := &TrelloItemId{
		boardId:     tg.options.BoardId,
		todoListId:  tg.options.TodoListId,
		doingListId: tg.options.DoingListId,
		doneListId:  tg.options.DoneListId,
	}

	if err := tg.AddTrelloIdSecret(trelloItemId); err != nil {
		return nil, err
	}

	log.Success("Adding secret keys for trello succeeded.")

	return buildState(tg), nil
}
