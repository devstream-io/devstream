package pluginengine

import (
	"errors"
	"os"

	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Destroy(continueDirectly bool) error {
	smgr, err := statemanager.NewManager()
	if err != nil {
		log.Debugf("Failed to get the manager: %s.", err)
		return err
	}

	changes, err := GetChangesForDestroy(smgr)
	if err != nil {
		log.Debugf("Get changes failed: %s.", err)
		return err
	}
	if len(changes) == 0 {
		log.Info("No tools have been deployed now. There is nothing to do.")
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
		return errors.New("some error(s) occurred during plugins destroy process")
	}

	log.Success("All plugins destroyed successfully.")
	return nil
}
