package main

import (
	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/internal/pkg/plugin/trellogithub"
)

// NAME is the name of this DevStream plugin.
const NAME = "trellogithub"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the installation of some trello-github-integ workflows.
func (p Plugin) Create(options *map[string]interface{}) (map[string]interface{}, error) {
	return trellogithub.Create(options)
}

// Update implements the installation of some trello-github-integ workflows.
func (p Plugin) Update(options *map[string]interface{}) (map[string]interface{}, error) {
	return trellogithub.Update(options)
}

// Read implements the healthy check of trello-github-integ workflows.
func (p Plugin) Read(options *map[string]interface{}) (map[string]interface{}, error) {
	return trellogithub.Read(options)
}

// Delete implements the installation of some trello-github-integ workflows.
func (p Plugin) Delete(options *map[string]interface{}) (bool, error) {
	return trellogithub.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
