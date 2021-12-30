package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/internal/pkg/backend"
	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/planmanager"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/internal/pkg/statemanager"
)

var installCMD = &cobra.Command{
	Use:   "install",
	Short: "Install tools defined in config file",
	Long:  `Install tools defined in config file`,
	Run:   installCMDFunc,
}

func installCMDFunc(cmd *cobra.Command, args []string) {
	conf := configloader.LoadConf(configFile)

	// init before installation
	err := pluginmanager.DownloadPlugins(conf)
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

	p := planmanager.NewPlan(smgr, conf)
	if len(p.Changes) == 0 {
		log.Println("it is nothing to do here")
		return
	}

	errsMap := p.Execute()
	if len(errsMap) == 0 {
		log.Println("=== all plugins Install/Uninstall/Reinstall process succeeded ===")
		log.Println("=== END ===")
		return
	}

	log.Println("=== some errors occurred during plugins Install/Uninstall/Reinstall process ===")
	for k, err := range errsMap {
		log.Printf("%s -> %s", k, err)
	}
	log.Println("=== END ===")
}
