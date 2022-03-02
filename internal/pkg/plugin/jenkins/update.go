package jenkins

import (
	"github.com/merico-dev/stream/pkg/util/log"
)

// Update updates jenkins with provided options.
func Update(options map[string]interface{}) (map[string]interface{}, error) {
	_, err := Delete(options)
	if err != nil {
		log.Errorf("Failed to delete the Jenkins: %s.", err)
		return nil, err
	}
	return Create(options)
}
