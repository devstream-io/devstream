package trello

import (
	"github.com/merico-dev/stream/pkg/util/log"
)

// Create creates Tello board and lists(todo/doing/done).
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts *Options
	var err error

	if opts, err = convertMap2Options(options); err != nil {
		return nil, err
	}

	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	trelloIds, err := CreateTrelloBoard(opts)
	if err != nil {
		return nil, err
	}
	log.Success("Creating trello board succeeded.")

	return buildState(opts, trelloIds), nil
}
