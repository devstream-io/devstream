package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/gitlabci/golang"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "gitlabci-golang"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create creates a GitLab CI workflow for Golang.
func (p Plugin) Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return golang.Create(options)
}

// Update updates the GitLab CI workflow for Golang.
func (p Plugin) Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return golang.Update(options)
}

// Read gets the state of the GitLab CI workflow for Golang.
func (p Plugin) Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return golang.Read(options)
}

// Delete deletes the GitLab CI workflow.
func (p Plugin) Delete(options configmanager.RawOptions) (bool, error) {
	return golang.Delete(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
