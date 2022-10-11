package pluginengine

import (
	"errors"
	"fmt"
	"os"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/interact"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Destroy(configFile string, continueDirectly bool, isForceDestroy bool) error {
	cfg, err := configmanager.NewManager(configFile).LoadConfig()
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("failed to load the config file")
	}

	if err := file.SetPluginDir(cfg.PluginDir); err != nil {
		log.Errorf("Error: %s.", err)
	}

	smgr, err := statemanager.NewManager(*cfg.State)
	if err != nil {
		log.Debugf("Failed to get the manager: %s.", err)
		return err
	}

	changes, err := GetChangesForDestroy(smgr, isForceDestroy)
	if err != nil {
		log.Debugf("Get changes failed: %s.", err)
		return err
	}
	if len(changes) == 0 {
		log.Info("No tools have been deployed now. There is nothing to do.")
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

	errsMap := execute(smgr, changes, true)
	if len(errsMap) != 0 {
		for k, e := range errsMap {
			log.Infof("%s -> %s", k, e)
		}
		return errors.New("some error(s) occurred during plugins destroy process")
	}

	log.Success("All plugins destroyed successfully.")
	return nil
}
