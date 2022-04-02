package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabci/golang"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "gitlabci-golang"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create creates a GitLab CI workflow for Golang.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return golang.Create(options)
}

// Update updates the GitLab CI workflow for Golang.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return golang.Update(options)
}

// Read gets the state of the GitLab CI workflow for Golang.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return golang.Read(options)
}

// Delete deletes the GitLab CI workflow.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return golang.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
