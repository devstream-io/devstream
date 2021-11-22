package main

import (
	"fmt"
	"log"

	"github.com/merico-dev/stream/internal/pkg/argocdapp"
)

// NAME is the name of this DevStream plugin.
const NAME = "argocdapps"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of an ArgoCD app.
func (p Plugin) Install(options *map[string]interface{}) {
	argocdapp.Install(options)
	log.Println("argocdapps install finished")
}

// Reinstall implements the installation of an ArgoCD app.
func (p Plugin) Reinstall(options *map[string]interface{}) {
	log.Println("mock: argocdapps reinstall finished")
}

// Uninstall Uninstall the installation of an ArgoCD app.
func (p Plugin) Uninstall(options *map[string]interface{}) {
	log.Println("mock: argocdapps uninstall finished")
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	fmt.Println("This is a plugin for DevStream. Use it with DevStream.")
}
