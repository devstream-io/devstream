package trello

import (
	"github.com/merico-dev/stream/pkg/util/log"
)

// Create creates Tello board and lists(todo/doing/done).
func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opt *Options
	var err error

	if opt, err = convertMap2Options(options); err != nil {
		return nil, err
	}

	if err := validateOptions(opt); err != nil {
		return nil, err
	}

	trelloIds, err := CreateTrelloBoard(opt)
	if err != nil {
		return nil, err
	}
	log.Success("Creating trello board succeeded.")

	return buildState(opt, trelloIds), nil
}
