package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/zentao"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "zentao"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of zentao.
func (p Plugin) Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return zentao.Create(options)
}

// Update implements the update of zentao.
func (p Plugin) Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return zentao.Update(options)
}

// Delete implements the delete of zentao.
func (p Plugin) Delete(options configmanager.RawOptions) (bool, error) {
	return zentao.Delete(options)
}

// Read implements the read of zentao.
func (p Plugin) Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return zentao.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
