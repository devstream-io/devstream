package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "jenkins"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of the jenkins.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return jenkins.Create(options)
}

// Update implements the update of the jenkins.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return jenkins.Update(options)
}

// Read implements read of the jenkins.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return jenkins.Read(options)
}

// Delete implements the delete of the jenkins.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return jenkins.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
