package main

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/plugin/argocdapp"
)

// NAME is the name of this DevStream plugin.
const NAME = "argocdapps"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of an ArgoCD app.
func (p Plugin) Create(options *map[string]interface{}) (bool, error) {
	return argocdapp.Install(options)
}

// Reinstall implements the installation of an ArgoCD app.
func (p Plugin) Update(options *map[string]interface{}) (bool, error) {
	return argocdapp.Reinstall(options)
}

// IsHealthy implements the healthy check of ArgoCD app.
func (p Plugin) Read(options *map[string]interface{}) (bool, error) {
	return argocdapp.IsHealthy(options)
}

// Delete Deletes the installation of an ArgoCD app.
func (p Plugin) Delete(options *map[string]interface{}) (bool, error) {
	return argocdapp.Uninstall(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
