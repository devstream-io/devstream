package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/devlake"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "devlake"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the installation of DevLake.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return devlake.Create(options)
}

// Update implements the installation of DevLake.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return devlake.Update(options)
}

// Read implements the healthy check of DevLake.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return devlake.Read(options)
}

// Delete Uninstall the installation of DevLake.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return devlake.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
