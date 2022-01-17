package pluginengine

import (
	"errors"
	"log"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/planmanager"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func Delete(fname string) error {
	cfg := configloader.LoadConf(fname)

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

	p := planmanager.NewDeletePlan(smgr, cfg)
	if len(p.Changes) == 0 {
		log.Println("Nothing needs to be deleted. There is nothing to do.")
		return nil
	}

	errsMap := execute(p)
	if len(errsMap) != 0 {
		err := errors.New("some error(s) occurred during plugins delete process")
		for k, e := range errsMap {
			log.Printf("%s -> %s", k, e)
		}
		return err
	}

	log.Println("All plugins deleted successfully.")
	return nil
}
