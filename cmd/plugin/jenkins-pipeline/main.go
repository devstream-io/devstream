package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/jenkinspipeline"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "jenkins-pipeline"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of jenkins-pipeline.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return jenkinspipeline.Create(options)
}

// Update implements the update of jenkins-pipeline.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return jenkinspipeline.Update(options)
}

// Delete implements the delete of jenkins-pipeline.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return jenkinspipeline.Delete(options)
}

// Read implements the read of jenkins-pipeline.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return jenkinspipeline.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
