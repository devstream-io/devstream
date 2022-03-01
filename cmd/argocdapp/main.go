package main

import (
	"github.com/merico-dev/stream/internal/pkg/plugin/argocdapp"
	"github.com/merico-dev/stream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "argocdapps"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the installation of an ArgoCD app.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return argocdapp.Create(options)
}

// Update implements the installation of an ArgoCD app.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return argocdapp.Update(options)
}

// Read implements the healthy check of ArgoCD app.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return argocdapp.Read(options)
}

// Delete Deletes the installation of an ArgoCD app.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return argocdapp.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
