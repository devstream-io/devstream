package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabcedocker"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "gitlab-ce-docker"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of gitlab-ce-docker.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return gitlabcedocker.Create(options)
}

// Update implements the update of gitlab-ce-docker.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return gitlabcedocker.Update(options)
}

// Delete implements the delete of gitlab-ce-docker.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return gitlabcedocker.Delete(options)
}

// Read implements the read of gitlab-ce-docker.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return gitlabcedocker.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
