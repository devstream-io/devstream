package main

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/plugin/kubeprometheus"
)

// NAME is the name of this DevStream plugin.
const NAME = "kube-prometheus"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of the kube-prometheus.
func (p Plugin) Install(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.Install(options)
}

// Reinstall implements the reinstallation of the kube-prometheus.
func (p Plugin) Reinstall(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.Reinstall(options)
}

// Uninstall implements the uninstallation of the kube-prometheus.
func (p Plugin) Uninstall(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.Uninstall(options)
}

// IsHealthy implements the healthy check of the kube-prometheus.
func (p Plugin) IsHealthy(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.IsHealthy(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
