package trello

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/pkg/util/log"
	"github.com/merico-dev/stream/pkg/util/trello"
)

func buildState(options *Options, ti *TrelloItemId) map[string]interface{} {
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

func buildReadState(options *Options) (map[string]interface{}, error) {
	c, err := trello.NewClient()
	if err != nil {
		return nil, err
	}
	listIds, err := c.GetBoardIdAndListId(options.Owner, options.Repo, options.KanbanBoardName)
	if err != nil {
		return nil, err
	}

	return listIds, nil
}

func convertMap2Options(options map[string]interface{}) (*Options, error) {
	var opt Options
	err := mapstructure.Decode(options, &opt)
	if err != nil {
		return nil, err
	}
	return &opt, nil
}

func validateOptions(options *Options) error {
	if errs := validate(options); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s.", e)
		}
		return fmt.Errorf("params are illegal")
	}
	return nil
}
