package trello

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/trello"
)

func getState(rawOptions configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := newOptions(rawOptions)
	if err != nil {
		return nil, err
	}

	c, err := trello.NewClient()
	if err != nil {
		return nil, err
	}
	idData, err := opts.Board.get(c)
	if err != nil {
		return nil, err
	}
	idDataMap, err := mapz.DecodeStructToMap(idData)
	if err != nil {
		return nil, err
	}

	// get board status
	resStatus := statemanager.ResourceStatus(idDataMap)
	output := make(statemanager.ResourceOutputs)
	output["boardId"] = fmt.Sprint(idData.boardID)
	output["todoListId"] = fmt.Sprint(idData.todoListID)
	output["doingListId"] = fmt.Sprint(idData.doingListID)
	output["doneListId"] = fmt.Sprint(idData.doneListID)
	resStatus.SetOutputs(output)

	// get ci status
	ciFileStatus, err := cifile.GetCIFileStatus(rawOptions)
	if err != nil {
		return nil, err
	}
	resStatus["ci"] = ciFileStatus
	return resStatus, nil
}
