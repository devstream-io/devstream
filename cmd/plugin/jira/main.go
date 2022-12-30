package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/jira"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "jira"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the installation of some jira-github-integ workflows.
func (p Plugin) Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return jira.Create(options)
}

// Update implements the installation of some jira-github-integ workflows.
func (p Plugin) Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return jira.Update(options)
}

// Read implements the healthy check of jira-github-integ workflows.
func (p Plugin) Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return jira.Read(options)
}

// Delete implements the installation of some jira-github-integ workflows.
func (p Plugin) Delete(options configmanager.RawOptions) (bool, error) {
	return jira.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
