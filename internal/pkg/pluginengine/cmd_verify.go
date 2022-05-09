package pluginengine

import (
	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/pluginmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Verify returns true if all the comments in this function are met
func Verify(configFile string) bool {
	// 1. loading config file succeeded
	cfg, err := configloader.LoadConf(configFile)
	if err != nil {
		log.Errorf("verify failed, error: %s", err)
	}

	if cfg == nil {
		return false
	}

	// 2. according to the config, all needed plugins exist
	err = pluginmanager.CheckLocalPlugins(cfg)
	if err != nil {
		log.Info(err)
		log.Info("Maybe you forgot to run \"dtm init\" first?")
		return false
	}
	// 3. can successfully create the state
	smgr, err := statemanager.NewManager()
	if err != nil {
		log.Errorf("Something is wrong with the state: %s.", err)
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
