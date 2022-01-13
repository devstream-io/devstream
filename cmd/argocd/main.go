package main

import (
	"log"

	"github.com/merico-dev/stream/internal/pkg/plugin/argocd"
)

// NAME is the name of this DevStream plugin.
const NAME = "argocd"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of ArgoCD.
func (p Plugin) Install(options *map[string]interface{}) (bool, error) {
	return argocd.Install(options)
}

// Reinstall implements the reinstallation of ArgoCD.
func (p Plugin) Reinstall(options *map[string]interface{}) (bool, error) {
	return argocd.Reinstall(options)
}

// Uninstall implements the uninstallation of ArgoCD.
func (p Plugin) Uninstall(options *map[string]interface{}) (bool, error) {
	return argocd.Uninstall(options)
}

// IsHealthy implements the healthy check of ArgoCD.
func (p Plugin) IsHealthy(options *map[string]interface{}) (bool, error) {
	return argocd.IsHealthy(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Printf("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
