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
func (p Plugin) Create(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.Install(options)
}

// Reinstall implements the reinstallation of the kube-prometheus.
func (p Plugin) Update(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.Reinstall(options)
}

// IsHealthy implements the healthy check of the kube-prometheus.
func (p Plugin) Read(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.IsHealthy(options)
}

// Uninstall implements the uninstallation of the kube-prometheus.
func (p Plugin) Delete(options *map[string]interface{}) (bool, error) {
	return kubeprometheus.Uninstall(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
