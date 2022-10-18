package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/kubeprometheus"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "kube-prometheus"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of the kube-prometheus.
func (p Plugin) Create(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return kubeprometheus.Create(options)
}

// Update implements the update of the kube-prometheus.
func (p Plugin) Update(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return kubeprometheus.Update(options)
}

// Read implements read of the kube-prometheus.
func (p Plugin) Read(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return kubeprometheus.Read(options)
}

// Delete implements the delete of the kube-prometheus.
func (p Plugin) Delete(options configmanager.RawOption) (bool, error) {
	return kubeprometheus.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
