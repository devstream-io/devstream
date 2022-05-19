package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/trello"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "trello"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the creation of trello board.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return trello.Create(options)
}

// Update implements the creation of trello board.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return trello.Update(options)
}

// Read implements the healthy check of trello board.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return trello.Read(options)
}

// Delete implements the creation of trello board.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return trello.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
