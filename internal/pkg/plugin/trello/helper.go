package trello

import (
	"fmt"

	"github.com/merico-dev/stream/pkg/util/trello"
)

func buildState(opts *Options, ti *TrelloItemId) map[string]interface{} {
	res := make(map[string]interface{})
	res["boardId"] = ti.boardId
	res["todoListId"] = ti.todoListId
	res["doingListId"] = ti.doingListId
	res["doneListId"] = ti.doneListId

	output := make(map[string]interface{})
	output["boardId"] = ti.boardId
	output["todoListId"] = ti.todoListId
	output["doingListId"] = ti.doingListId
	output["doneListId"] = ti.doneListId

	res["outputs"] = output

	return res
}

func buildReadState(opts *Options) (map[string]interface{}, error) {
	c, err := trello.NewClient()
	if err != nil {
		return nil, err
	}
	listIds, err := c.GetBoardIdAndListId(opts.Owner, opts.Repo, opts.KanbanBoardName)
	if err != nil {
		return nil, err
	}

	output := make(map[string]interface{})
	output["boardId"] = fmt.Sprint(listIds["boardId"])
	output["todoListId"] = fmt.Sprint(listIds["todoListId"])
	output["doingListId"] = fmt.Sprint(listIds["doingListId"])
	output["doneListId"] = fmt.Sprint(listIds["doneListId"])

	listIds["outputs"] = output
	return listIds, nil
}
