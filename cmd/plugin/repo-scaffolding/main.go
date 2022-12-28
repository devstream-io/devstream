package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/reposcaffolding"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "repo-scaffolding"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of repo-scaffolding.
func (p Plugin) Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return reposcaffolding.Create(options)
}

// Update implements the update of repo-scaffolding.
func (p Plugin) Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return reposcaffolding.Update(options)
}

// Delete implements the delete of repo-scaffolding.
func (p Plugin) Delete(options configmanager.RawOptions) (bool, error) {
	return reposcaffolding.Delete(options)
}

// Read implements the read of repo-scaffolding.
func (p Plugin) Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return reposcaffolding.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
