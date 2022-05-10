package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/helmgeneric"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "helm-generic"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of helm-generic.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return helmgeneric.Create(options)
}

// Update implements the update of helm-generic.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return helmgeneric.Update(options)
}

// Delete implements the delete of helm-generic.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return helmgeneric.Delete(options)
}

// Read implements the read of helm-generic.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return helmgeneric.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
