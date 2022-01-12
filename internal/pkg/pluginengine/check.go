package pluginengine

import (
	"log"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

// CheckHealthy returns true while all tools are healthy
func CheckHealthy(fname string) bool {
	cfg := configloader.LoadConf(fname)
	allHealthy := true

	for _, tool := range cfg.Tools {
		healthy, err := IsHealthy(&tool)
		if err != nil {
			allHealthy = false
			log.Printf("failed to check healthy for the tool: %s, got error: %s", tool.Name, err)
			continue
		}
		if healthy {
			log.Printf("the tool %s is healthy", tool.Name)
			continue
		}
		allHealthy = false
		log.Printf("the tool %s is not healthy", tool.Name)
	}

	return allHealthy
}
