package pluginengine

import (
	"errors"
	"fmt"
	"os"

	"github.com/merico-dev/stream/internal/pkg/log"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func Apply(configFile string, continueDirectly bool) error {
	cfg := configloader.LoadConf(configFile)
	if cfg == nil {
		return fmt.Errorf("failed to load the config file")
	}

	err := pluginmanager.CheckLocalPlugins(cfg)
	if err != nil {
		log.Errorf("Error checking required plugins. Maybe you forgot to run \"dtm init\" first?")
		return err
	}

	// use default local backend for now.
	b, err := backend.GetBackend(backend.BackendLocal)
	if err != nil {
		return err
	}
	// create a state manager using the default local backend
	smgr := statemanager.NewManager(b)

	changes, err := GetChangesForApply(smgr, cfg)
	if err != nil {
		log.Debugf("Get changes for apply failed: %s", err)
		return err
	}
	if len(changes) == 0 {
		log.Info("No changes done since last apply. There is nothing to do.")
		return nil
	}

	if !continueDirectly {
		userInput := readUserInput()
		if userInput == "n" {
			os.Exit(0)
		}
	}

	errsMap := execute(smgr, changes)
	if len(errsMap) != 0 {
		for k, e := range errsMap {
			log.Infof("%s -> %s", k, e)
		}
		return errors.New("some error(s) occurred during plugins apply process")
	}

	log.Success("All plugins applied successfully.")
	return nil
}
