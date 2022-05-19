package pluginengine

import (
	"errors"
	"fmt"
	"os"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/pluginmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Apply(configFile string, continueDirectly bool) error {
	cfg, err := configloader.LoadConf(configFile)
	if err != nil {
		return err
	}

	if cfg == nil {
		return fmt.Errorf("failed to load the config file")
	}

	err = pluginmanager.CheckLocalPlugins(cfg)
	if err != nil {
		log.Error("Error checking required plugins. Maybe you forgot to run \"dtm init\" first?")
		return err
	}

	smgr, err := statemanager.NewManager(*cfg.State)
	if err != nil {
		log.Debugf("Failed to get the manager: %s.", err)
		return err
	}

	changes, err := GetChangesForApply(smgr, cfg)
	if err != nil {
		log.Debugf("Get changes for apply failed: %s.", err)
		return err
	}
	if len(changes) == 0 {
		log.Info("No changes done since last apply. There is nothing to do.")
		return nil
	}

	for _, change := range changes {
		log.Info(change.Description)
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
			log.Errorf("Errors Map: key(%s) -> value(%s)", k, e)
		}
		return errors.New("some error(s) occurred during plugins apply process")
	}

	log.Success("All plugins applied successfully.")
	return nil
}
