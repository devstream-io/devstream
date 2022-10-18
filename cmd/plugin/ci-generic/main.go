package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/cigeneric"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "ci-generic"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of ci-generic.
func (p Plugin) Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return cigeneric.Create(options)
}

// Update implements the update of ci-generic.
func (p Plugin) Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return cigeneric.Update(options)
}

// Delete implements the delete of ci-generic.
func (p Plugin) Delete(options configmanager.RawOptions) (bool, error) {
	return cigeneric.Delete(options)
}

// Read implements the read of ci-generic.
func (p Plugin) Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return cigeneric.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
