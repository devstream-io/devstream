package pluginengine

import (
	"errors"
	"fmt"
	"log"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/planmanager"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func Apply(fname string) error {
	cfg := configloader.LoadConf(fname)
	if cfg == nil {
		return fmt.Errorf("failed to load the config file")
	}

	err := pluginmanager.CheckLocalPlugins(cfg)
	if err != nil {
		log.Printf("Error checking required plugins. Maybe you forgot to run \"dtm init\" first?")
		return err
	}

	// use default local backend for now.
	b, err := backend.GetBackend("local")
	if err != nil {
		return err
	}

	smgr := statemanager.NewManager(b)

	p := planmanager.NewPlan(smgr, cfg)
	if len(p.Changes) == 0 {
		log.Println("No changes done since last apply. There is nothing to do.")
		return nil
	}

	errsMap := execute(p)
	if len(errsMap) != 0 {
		for k, e := range errsMap {
			log.Printf("%s -> %s", k, e)
		}
		return errors.New("some error(s) occurred during plugins delete process")
	}

	log.Println("All plugins applied successfully.")
	return nil
}
