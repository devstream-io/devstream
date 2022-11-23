package pluginengine

import (
	"errors"
	"os"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/pluginmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/interact"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const askUserIfContinue string = "Continue? [y/n]"

func Apply(configFile string, continueDirectly bool) error {
	cfg, err := configmanager.NewManager(configFile).LoadConfig()
	if err != nil {
		return err
	}

	err = pluginmanager.CheckLocalPlugins(cfg.Tools)
	if err != nil {
		log.Error(`Error checking required plugins. Maybe you forgot to run "dtm init" first?`)
		return err
	}
	smgr, err := statemanager.NewManager(*cfg.Config.State)
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
		continued := interact.AskUserIfContinue(askUserIfContinue)
		if !continued {
			os.Exit(0)
		}
	}

	errsMap := execute(smgr, changes, false)
	if len(errsMap) != 0 {
		for k, e := range errsMap {
			log.Errorf("Errors Map: key(%s) -> value(%s)", k, e)
		}
		return errors.New("some error(s) occurred during plugins apply process")
	}

	log.Success("All plugins applied successfully.")
	return nil
}
