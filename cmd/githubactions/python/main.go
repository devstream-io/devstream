package main

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/plugin/githubactions/python"
)

// NAME is the name of this DevStream plugin.
const NAME = "githubactions-python"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of some GitHub Actions workflows.
func (p Plugin) Create(options *map[string]interface{}) (bool, error) {
	return python.Install(options)
}

// Reinstall implements the installation of some GitHub Actions workflows.
func (p Plugin) Update(options *map[string]interface{}) (bool, error) {
	return python.Reinstall(options)
}

// IsHealthy implements the healthy check of GitHub Actions workflows.
func (p Plugin) Read(options *map[string]interface{}) (bool, error) {
	return python.IsHealthy(options)
}

// Uninstall implements the installation of some GitHub Actions workflows.
func (p Plugin) Delete(options *map[string]interface{}) (bool, error) {
	return python.Uninstall(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
