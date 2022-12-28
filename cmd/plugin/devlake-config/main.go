package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/devlakeconfig"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "devlake-config"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of devlake-config.
func (p Plugin) Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return devlakeconfig.Create(options)
}

// Update implements the update of devlake-config.
func (p Plugin) Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return devlakeconfig.Update(options)
}

// Delete implements the delete of devlake-config.
func (p Plugin) Delete(options configmanager.RawOptions) (bool, error) {
	return devlakeconfig.Delete(options)
}

// Read implements the read of devlake-config.
func (p Plugin) Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return devlakeconfig.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
