package pluginengine

import (
	"fmt"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
	"github.com/merico-dev/stream/pkg/util/log"
)

// Verify returns true while all tools are healthy
func Verify(configFile string) (bool, error) {
	cfg := configloader.LoadConf(configFile)
	if cfg == nil {
		return false, fmt.Errorf("failed to load the config file")
	}

	err := pluginmanager.CheckLocalPlugins(cfg)
	if err != nil {
		log.Errorf("Error checking required plugins. Maybe you forgot to run \"dtm init\" first?")
		return false, err
	}

	smgr, err := statemanager.NewManager()
	if err != nil {
		log.Debugf("Failed to get the manager: %s.", err)
		return false, err
	}

	changes, err := GetChangesForApply(smgr, cfg)
	if err != nil {
		log.Debugf("Get changes for apply failed: %s.", err)
		return false, err
	}
	if len(changes) == 0 {
		log.Info("All plugins is healthy now.")
		return true, nil
	}

	for _, c := range changes {
		log.Infof("The plugin < %s > has been changed, need to %s.", c.Tool.Name, c.ActionName)
	}
	return false, nil
}
