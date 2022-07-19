package helmgeneric

import "github.com/devstream-io/devstream/internal/pkg/plugininstaller"

// return empty
func getEmptyState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	retMap := make(map[string]interface{})
	return retMap, nil
}
