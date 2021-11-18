package main

import (
	"fmt"
	"log"

	"github.com/ironcore864/openstream/internal/pkg/githubactions"
)

const NAME = "githubactions"

type Plugin string

func (p Plugin) Install(options *map[string]interface{}) {
	githubactions.Install(options)
	log.Println("github actions install finished")
}

func (p Plugin) Reinstall(options *map[string]interface{}) {
	log.Println("mock: github actions reinstall finished")
}

func (p Plugin) Uninstall(options *map[string]interface{}) {
	log.Println("mock: github actions uninstall finished")
}

var OpenStreamPlugin Plugin

func main() {
	fmt.Println("This is a plugin for OpenStream. Use it with OpenStream.")
}
