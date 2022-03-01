package main

import (
	"github.com/merico-dev/stream/internal/pkg/plugin/argocd"
	"github.com/merico-dev/stream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "argocd"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of ArgoCD.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return argocd.Create(options)
}

// Update implements the update of ArgoCD.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return argocd.Update(options)
}

// Delete implements the delete of ArgoCD.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return argocd.Delete(options)
}

// Read implements the read of ArgoCD.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return argocd.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
