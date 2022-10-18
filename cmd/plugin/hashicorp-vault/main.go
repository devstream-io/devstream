package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/hashicorpvault"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "hashicorp-vault"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of hashicorp-vault.
func (p Plugin) Create(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return hashicorpvault.Create(options)
}

// Update implements the update of hashicorp-vault.
func (p Plugin) Update(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return hashicorpvault.Update(options)
}

// Delete implements the delete of hashicorp-vault.
func (p Plugin) Delete(options configmanager.RawOption) (bool, error) {
	return hashicorpvault.Delete(options)
}

// Read implements the read of hashicorp-vault.
func (p Plugin) Read(options configmanager.RawOption) (statemanager.ResourceStatus, error) {
	return hashicorpvault.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
