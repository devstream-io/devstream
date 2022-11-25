package trello

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	. "github.com/devstream-io/devstream/internal/pkg/plugin/common"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Delete delete trello board and lists
func Delete(options configmanager.RawOptions) (bool, error) {
	var err error
	defer func() {
		HandleErrLogsWithPlugin(err, Name)
	}()

	var opts Options
	if err = mapstructure.Decode(options, &opts); err != nil {
		return false, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return false, fmt.Errorf("opts are illegal")
	}

	if err := DeleteTrelloBoard(&opts); err != nil {
		return false, err
	}
	return true, nil
}
