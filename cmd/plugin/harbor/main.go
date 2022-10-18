package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/harbor"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "harbor"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of harbor.
func (p Plugin) Create(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return harbor.Create(options)
}

// Update implements the update of harbor.
func (p Plugin) Update(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return harbor.Update(options)
}

// Delete implements the delete of harbor.
func (p Plugin) Delete(options configmanager.RawOption) (bool, error) {
	return harbor.Delete(options)
}

// Read implements the read of harbor.
func (p Plugin) Read(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return harbor.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
