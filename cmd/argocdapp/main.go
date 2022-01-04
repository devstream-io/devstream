package main

import (
	"log"

	"github.com/merico-dev/stream/internal/pkg/plugin/argocdapp"
)

// NAME is the name of this DevStream plugin.
const NAME = "argocdapps"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of an ArgoCD app.
func (p Plugin) Install(options *map[string]interface{}) (bool, error) {
	return argocdapp.Install(options)
}

// Reinstall implements the installation of an ArgoCD app.
func (p Plugin) Reinstall(options *map[string]interface{}) (bool, error) {
	return argocdapp.Reinstall(options)
}

// Uninstall Uninstall the installation of an ArgoCD app.
func (p Plugin) Uninstall(options *map[string]interface{}) (bool, error) {
	return argocdapp.Uninstall(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Printf("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
