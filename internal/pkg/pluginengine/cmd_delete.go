package pluginengine

import (
	"errors"
	"fmt"
	"os"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func Remove(configFile string, continueDirectly bool) error {
	cfg := configloader.LoadConf(configFile)
	if cfg == nil {
		return fmt.Errorf("failed to load the config file")
	}

	err := pluginmanager.CheckLocalPlugins(cfg)
	if err != nil {
		log.Errorf("Error checking required plugins. Maybe you forgot to run \"dtm init\" first?")
		return err
	}

	smgr, err := statemanager.NewManager()
	if err != nil {
		log.Debugf("Failed to get the manager: %s.", err)
		return err
	}

	changes, err := GetChangesForDelete(smgr, cfg)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		log.Info("Nothing needs to be deleted. There is nothing to do.")
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
		err := errors.New("some error(s) occurred during plugins delete process")
		for k, e := range errsMap {
			log.Infof("%s -> %s.", k, e)
		}
		return err
	}

	log.Success("All plugins deleted successfully.")
	return nil
}
