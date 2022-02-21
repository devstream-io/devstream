package pluginengine

import (
	"fmt"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

// Verify returns true while all tools are healthy
func Verify(fname string) (bool, error) {
	cfg := configloader.LoadConf(fname)
	if cfg == nil {
		return false, fmt.Errorf("failed to load the config file")
	}

	allHealthy := true
	// TODO(daniel-hutao): waiting for refactor with verify
	//for _, tool := range cfg.Tools {
	//	healthy, err := Read(&tool)
	//	if err != nil {
	//		allHealthy = false
	//		log.Errorf("failed to check healthy for the tool: %s, got error: %s", tool.Name, err)
	//		continue
	//	}
	//	if healthy {
	//		log.Successf("the tool %s is healthy", tool.Name)
	//		continue
	//	}
	//	allHealthy = false
	//	log.Warnf("the tool %s is not healthy", tool.Name)
	//}

	return allHealthy, nil
}
