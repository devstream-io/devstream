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

func Run(fname string) error {
	cfg := configloader.LoadConf(fname)

	// init before installation
	err := pluginmanager.DownloadPlugins(cfg)
	if err != nil {
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
		err := errors.New("some error(s) occurred during plugins apply process")
		for k, e := range errsMap {
			log.Printf("%s -> %s", k, e)
		}
		return err
	}

	log.Println("All plugins applied successfully.")
	return nil
}
