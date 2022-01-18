package main

import (
	"log"

	"github.com/merico-dev/stream/internal/pkg/plugin/kubeprometheus"
)

// NAME is the name of this DevStream plugin.
const NAME = "kube-prometheus"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of some GitHub Actions workflows.
func (p Plugin) Install(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.Install(options)
}

// Reinstall implements the reinstallation of some GitHub Actions workflows.
func (p Plugin) Reinstall(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.Reinstall(options)
}

// Uninstall implements the uninstallation of some GitHub Actions workflows.
func (p Plugin) Uninstall(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.Uninstall(options)
}

// IsHealthy implements the healthy check of GitHub Actions workflows.
func (p Plugin) IsHealthy(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.IsHealthy(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Printf("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
