package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/jenkinspipeline"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "jenkins-pipeline"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of jenkins-pipeline.
func (p Plugin) Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return jenkinspipeline.Create(options)
}

// Update implements the update of jenkins-pipeline.
func (p Plugin) Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return jenkinspipeline.Update(options)
}

// Delete implements the delete of jenkins-pipeline.
func (p Plugin) Delete(options configmanager.RawOptions) (bool, error) {
	return jenkinspipeline.Delete(options)
}

// Read implements the read of jenkins-pipeline.
func (p Plugin) Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return jenkinspipeline.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
