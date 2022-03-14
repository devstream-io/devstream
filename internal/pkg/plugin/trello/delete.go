package trello

import "github.com/merico-dev/stream/pkg/util/log"

// Delete does not remove trello board and lists
func Delete(options map[string]interface{}) (bool, error) {
	log.Info("Delelte will do nothing, because someone might be using the existing board.")
	return true, nil
}
