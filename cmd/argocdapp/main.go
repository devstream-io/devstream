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
	log.Printf("%s install finished", NAME)
}

// Reinstall implements the installation of an ArgoCD app.
func (p Plugin) Reinstall(options *map[string]interface{}) {
	log.Printf("mock: %s reinstall finished", NAME)
}

// Uninstall Uninstall the installation of an ArgoCD app.
func (p Plugin) Uninstall(options *map[string]interface{}) {
	log.Printf("mock: %s uninstall finished", NAME)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	fmt.Printf("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
