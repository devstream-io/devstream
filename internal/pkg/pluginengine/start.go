package pluginengine

import (
	"fmt"
	"log"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/planmanager"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

type CmdType string

const (
	Init  CmdType = "init"
	Apply CmdType = "apply"
)

func Do(doType CmdType) error {
	cfg, err := configloader.LoadConfig()
	if err != nil {
		return err
	}

	switch doType {
	case Init:
		return initz(cfg)
	case Apply:
		return apply(cfg)
	default:
		return fmt.Errorf("illegal CmdType")
	}
}

// initz is same to init
func initz(cfg *configloader.Config) error {
	err := pluginmanager.DownloadPlugins(cfg)
	if err != nil {
		return err
	}
	log.Println("=== initialize finished ===")
	return nil
}

func apply(cfg *configloader.Config) error {
	// initialize before installation
	if err := initz(cfg); err != nil {
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
		log.Println("it is nothing to do here")
		return nil
	}

	errsMap := p.Execute()
	if len(errsMap) == 0 {
		log.Println("=== all plugins Install/Uninstall/Reinstall process succeeded ===")
		log.Println("=== END ===")
		return nil
	}

	log.Println("=== some errors occurred during plugins Install/Uninstall/Reinstall process ===")
	for k, err := range errsMap {
		log.Printf("%s -> %s", k, err)
	}
	log.Println("=== END ===")
	return nil
}
