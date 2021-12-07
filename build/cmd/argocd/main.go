package main

import (
	"fmt"
	"log"

	"github.com/merico-dev/stream/internal/pkg/argocd"
)

// NAME is the name of this DevStream plugin.
const NAME = "argocd"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of ArgoCD.
func (p Plugin) Install(options *map[string]interface{}) {
	argocd.Install(options)
	log.Printf("%s install finished", NAME)
}

// Reinstall implements the reinstallation of ArgoCD.
func (p Plugin) Reinstall(options *map[string]interface{}) {
	log.Printf("mock: %s reinstall finished", NAME)
}

// Uninstall implements the uninstallation of ArgoCD.
func (p Plugin) Uninstall(options *map[string]interface{}) {
	log.Printf("mock: %s uninstall finished", NAME)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	fmt.Printf("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
