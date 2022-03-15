package trello

import (
	"fmt"

	"github.com/merico-dev/stream/pkg/util/log"
	"github.com/merico-dev/stream/pkg/util/trello"
)

// Options is the struct for configurations of the trellogithub plugin.
type Options struct {
	Owner           string
	Repo            string
	KanbanBoardName string
}

type TrelloItemId struct {
	boardId     string
	todoListId  string
	doingListId string
	doneListId  string
}

// CreateTrelloItems create board, and lists will be created automatically
func CreateTrelloBoard(options *Options) (*TrelloItemId, error) {
	c, err := trello.NewClient()
	if err != nil {
		return nil, err
	}

	board, err := c.CreateBoard(options.KanbanBoardName, options.Owner, options.Repo)
	if err != nil {
		return nil, err
	}

	lists, err := board.GetLists()
	if err != nil {
		return nil, err
	}
	if len(lists) != 3 {
		log.Errorf("Unknown lists format: len==%d.", len(lists))
		return nil, fmt.Errorf("unknown lists format: len==%d", len(lists))
	}

	// (NOTICE) We need to get the list ids and then continue the following operations.
	// By default, creating a board will come with 3 lists, and they are ordered,
	// so we can only get it in this way for now.
	// It should be noted that if the implementation of trello changes one day,
	// some errors may occurred here, and we need to pay attention to this.
	todo := lists[0].ID
	doing := lists[1].ID
	done := lists[2].ID

	log.Debugf("Lists: To Do(%s), Doing(%s), Done(%s).", todo, doing, done)

	return &TrelloItemId{
		boardId:     board.ID,
		todoListId:  todo,
		doingListId: doing,
		doneListId:  done,
	}, nil
}
