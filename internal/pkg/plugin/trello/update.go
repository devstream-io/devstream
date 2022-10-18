package trello

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Update recreate trello board and lists.
func Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	if err := DeleteTrelloBoard(&opts); err != nil {
		return nil, err
	}

	trelloIds, err := CreateTrelloBoard(&opts)
	if err != nil {
		return nil, err
	}

	return buildStatus(&opts, trelloIds), nil
}
