package main

import (
	"log"

	"github.com/devstream-io/devstream/internal/pkg/plugin/reposcaffolding/github/golang"
)

// NAME is the name of this DevStream plugin.
const NAME = "github-repo-scaffolding-golang"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the installation of the github-repo-scaffolding-golang.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return golang.Create(options)
}

// Update implements the reinstallation of the github-repo-scaffolding-golang.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return golang.Update(options)
}

// Read implements the healthy check of the github-repo-scaffolding-golang.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return golang.Read(options)
}

// Delete implements the uninstallation of the github-repo-scaffolding-golang.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return golang.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Printf("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
