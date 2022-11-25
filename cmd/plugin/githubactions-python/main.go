package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/githubactions/python"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "githubactions-python"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of some githubactions-python.
func (p Plugin) Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return python.Create(options)
}

// Update implements the update of some githubactions-python.
func (p Plugin) Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return python.Update(options)
}

// Read implements the read of githubactions-python.
func (p Plugin) Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return python.Read(options)
}

// Delete implements the delete of some githubactions-python.
func (p Plugin) Delete(options configmanager.RawOptions) (bool, error) {
	return python.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
