package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/githubactions/python"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "githubactions-python"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the installation of some GitHub Actions workflows.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return python.Create(options)
}

// Update implements the installation of some GitHub Actions workflows.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return python.Update(options)
}

// Read implements the healthy check of GitHub Actions workflows.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return python.Read(options)
}

// Delete implements the installation of some GitHub Actions workflows.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return python.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
