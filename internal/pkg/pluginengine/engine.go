package pluginengine

import (
	"log"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/planmanager"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

func Run(fname string) {
	cfg := configloader.LoadConf(fname)

	// init before installation
	err := pluginmanager.DownloadPlugins(cfg)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	// use default local backend for now.
	b, err := backend.GetBackend("local")
	if err != nil {
		log.Fatal(err)
	}
	smgr := statemanager.NewManager(b)

	p := planmanager.NewPlan(smgr, cfg)
	if len(p.Changes) == 0 {
		log.Println("No changes done since last apply. There is nothing to do.")
		return
	}

	errsMap := execute(p)
	if len(errsMap) == 0 {
		log.Println("All plugins applied successfully.")
		return
	}

	log.Println("Some errors occurred during plugins apply process.")
	for k, err := range errsMap {
		log.Printf("%s -> %s", k, err)
	}
}
