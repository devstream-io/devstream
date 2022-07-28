package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabci/java"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "gitlabci-java"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of gitlabci-java.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return java.Create(options)
}

// Update implements the update of gitlabci-java.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return java.Update(options)
}

// Delete implements the delete of gitlabci-java.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return java.Delete(options)
}

// Read implements the read of gitlabci-java.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return java.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
