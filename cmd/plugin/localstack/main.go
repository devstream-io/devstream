package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/localstack"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "localstack"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of localstack.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return localstack.Create(options)
}

// Update implements the update of localstack.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return localstack.Update(options)
}

// Delete implements the delete of localstack.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return localstack.Delete(options)
}

// Read implements the read of localstack.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return localstack.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
