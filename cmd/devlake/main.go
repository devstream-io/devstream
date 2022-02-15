package main

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/plugin/devlake"
)

// NAME is the name of this DevStream plugin.
const NAME = "devlake"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of DevLake.
func (p Plugin) Create(options *map[string]interface{}) (bool, error) {
	return devlake.Install(options)
}

// Reinstall implements the installation of DevLake.
func (p Plugin) Update(options *map[string]interface{}) (bool, error) {
	return devlake.Reinstall(options)
}

// IsHealthy implements the healthy check of DevLake.
func (p Plugin) Read(options *map[string]interface{}) (bool, error) {
	return devlake.IsHealthy(options)
}

// Uninstall Uninstall the installation of DevLake.
func (p Plugin) Delete(options *map[string]interface{}) (bool, error) {
	return devlake.Uninstall(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
