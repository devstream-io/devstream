package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/tekton"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "tekton"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of tekton.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return tekton.Create(options)
}

// Update implements the update of tekton.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return tekton.Update(options)
}

// Delete implements the delete of tekton.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return tekton.Delete(options)
}

// Read implements the read of tekton.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return tekton.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
