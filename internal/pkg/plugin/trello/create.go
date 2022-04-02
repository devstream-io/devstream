package trello

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create creates Tello board and lists(todo/doing/done).
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options

	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	// validate parameters
	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	trelloIds, err := CreateTrelloBoard(&opts)
	if err != nil {
		return nil, err
	}
	log.Success("Creating trello board succeeded.")

	return buildState(&opts, trelloIds), nil
}
