package trello

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/trello"
)

// options is the struct for configurations of the trellogithub plugin.
type options struct {
	Scm   *git.RepoInfo `validate:"required" mapstructure:"scm"`
	Board *board        `validate:"required" mapstructure:"board"`

	// used in package
	CIFileConfig *cifile.CIFileConfig `mapstructure:"ci"`
}

type board struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
}

type boardIDInfo struct {
	boardID     string
	todoListID  string
	doingListID string
	doneListID  string
}

func newOptions(rawOptions configmanager.RawOptions) (*options, error) {
	var opts options
	if err := mapstructure.Decode(rawOptions, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func (b board) create(trelloClient trello.TrelloAPI) error {
	trelloBoard, err := trelloClient.Get(b.Name, b.Description)
	if err != nil {
		return err
	}
	if trelloBoard != nil {
		log.Debugf("Board already exists, description: %s, kanbanName: %s.", b.Description, b.Name)
		return nil
	}
	_, err = trelloClient.Create(b.Name, b.Description)
	return err
}

func (b board) delete(trelloClient trello.TrelloAPI) error {
	trelloBoard, err := trelloClient.Get(b.Name, b.Description)
	if err != nil {
		return err
	}
	if trelloBoard != nil {
		return trelloBoard.Delete()
	}
	return nil
}

func (b board) get(trelloClient trello.TrelloAPI) (*boardIDInfo, error) {
	const defaultTrelloListNum = 3
	trelloBoard, err := trelloClient.Get(b.Name, b.Description)
	if err != nil {
		return nil, err
	}
	boardList, err := trelloBoard.GetLists()
	if err != nil {
		return nil, err
	}
	if len(boardList) != defaultTrelloListNum {
		log.Errorf("Unknown lists format: len==%d.", len(boardList))
		return nil, fmt.Errorf("unknown lists format: len==%d", len(boardList))
	}
	idData := &boardIDInfo{
		boardID:     trelloBoard.ID,
		todoListID:  boardList[0].ID,
		doingListID: boardList[1].ID,
		doneListID:  boardList[2].ID,
	}
	return idData, nil
}
