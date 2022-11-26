package trellogithub

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Update remove and set up trello-github-integ workflows.
func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
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

	return buildStatus(tg), nil
}
