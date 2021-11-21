package main

import (
	"fmt"
	"log"

	"github.com/merico-dev/stream/internal/pkg/githubactions"
)

// NAME is the name of this DevStream plugin.
const NAME = "githubactions"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Install implements the installation of some GitHub Actions workflows.
func (p Plugin) Install(options *map[string]interface{}) {
	githubactions.Install(options)
	log.Println("github actions install finished")
}

// Reinstall implements the installation of some GitHub Actions workflows.
func (p Plugin) Reinstall(options *map[string]interface{}) {
	log.Println("mock: github actions reinstall finished")
}

// Uninstall implements the installation of some GitHub Actions workflows.
func (p Plugin) Uninstall(options *map[string]interface{}) {
	log.Println("mock: github actions uninstall finished")
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	fmt.Println("This is a plugin for DevStream. Use it with DevStream.")
}
