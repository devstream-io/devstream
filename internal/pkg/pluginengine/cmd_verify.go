package pluginengine

import (
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Verify returns true if all the comments in this function are met
func Verify(configFile string) bool {
	// 1. loading config file succeeded
	cfg := configloader.LoadConf(configFile)
	if cfg == nil {
		return false
	}

	// 2. according to the config, all needed plugins exist
	err := pluginmanager.CheckLocalPlugins(cfg)
	if err != nil {
		log.Info(err)
		log.Info("Maybe you forgot to run \"dtm init\" first?")
		return false
	}
	// 3. can successfully create the state
	smgr, err := statemanager.NewManager()
	if err != nil {
		log.Debugf("Failed to get the manager: %s.", err)
		return false
	}

	// 4. if the config, state and the resources/tools are exactly the same
	changes, err := GetChangesForApply(smgr, cfg)
	if err != nil {
		log.Errorf("Get changes failed: %s.", err)
		return false
	}
	if len(changes) != 0 {
		for _, change := range changes {
			log.Info(change.Description)
		}
		return false
	}
	return true
}
