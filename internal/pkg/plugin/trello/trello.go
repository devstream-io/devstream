package trello

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/trello"
)

// Options is the struct for configurations of the trellogithub plugin.
type Options struct {
	Owner           string `validate:"required_without=Org"`
	Org             string `validate:"required_without=Owner"`
	Repo            string `validate:"required"`
	KanbanBoardName string
}

type TrelloItemId struct {
	boardId     string
	todoListId  string
	doingListId string
	doneListId  string
}

// CreateTrelloBoard CreateTrelloItems create board, and lists will be created automatically
func CreateTrelloBoard(opts *Options) (*TrelloItemId, error) {
	c, err := trello.NewClient()
	if err != nil {
		return nil, err
	}

	exist, err := c.CheckBoardExists(opts.Owner, opts.Repo, opts.KanbanBoardName)
	if err != nil {
		return nil, err
	}

	if exist {
		log.Infof("Board already exists, owner: %s, repo: %s, kanbanName: %s.", opts.Owner, opts.Repo, opts.KanbanBoardName)
		listIds, err := c.GetBoardIdAndListId(opts.Owner, opts.Repo, opts.KanbanBoardName)
		if err != nil {
			return nil, err
		}
		return &TrelloItemId{
			boardId:     fmt.Sprint(listIds["boardId"]),
			todoListId:  fmt.Sprint(listIds["todoListId"]),
			doingListId: fmt.Sprint(listIds["doingListId"]),
			doneListId:  fmt.Sprint(listIds["doneListId"]),
		}, nil
	}

	board, err := c.CreateBoard(opts.KanbanBoardName, opts.Owner, opts.Repo)
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

// DeleteTrelloBoard delete specified board
func DeleteTrelloBoard(opts *Options) error {
	c, err := trello.NewClient()
	if err != nil {
		return err
	}
	return c.CheckAndDeleteBoard(opts.Owner, opts.Repo, opts.KanbanBoardName)
}
func buildStatus(opts *Options, ti *TrelloItemId) statemanager.ResourceStatus {
	resStatus := make(statemanager.ResourceStatus)
	resStatus["boardId"] = ti.boardId
	resStatus["todoListId"] = ti.todoListId
	resStatus["doingListId"] = ti.doingListId
	resStatus["doneListId"] = ti.doneListId

	outputs := make(statemanager.ResourceOutputs)
	outputs["boardId"] = ti.boardId
	outputs["todoListId"] = ti.todoListId
	outputs["doingListId"] = ti.doingListId
	outputs["doneListId"] = ti.doneListId

	resStatus.SetOutputs(outputs)

	return resStatus
}

func buildReadStatus(opts *Options) (statemanager.ResourceStatus, error) {
	c, err := trello.NewClient()
	if err != nil {
		return nil, err
	}
	listIds, err := c.GetBoardIdAndListId(opts.Owner, opts.Repo, opts.KanbanBoardName)
	if err != nil {
		return nil, err
	}

	resStatus := statemanager.ResourceStatus(listIds)
	output := make(statemanager.ResourceOutputs)
	output["boardId"] = fmt.Sprint(listIds["boardId"])
	output["todoListId"] = fmt.Sprint(listIds["todoListId"])
	output["doingListId"] = fmt.Sprint(listIds["doingListId"])
	output["doneListId"] = fmt.Sprint(listIds["doneListId"])

	resStatus.SetOutputs(output)
	return listIds, nil
}
