package main

import (
	"fmt"
	"log"

	"github.com/ironcore864/openstream/internal/pkg/argocdapp"
)

const NAME = "argocdapps"

type Plugin string

func (p Plugin) Install(options *map[string]interface{}) {
	argocdapp.Install(options)
	log.Println("argocdapps install finished")
}

func (p Plugin) Reinstall(options *map[string]interface{}) {
	log.Println("mock: argocdapps reinstall finished")
}

func (p Plugin) Uninstall(options *map[string]interface{}) {
	log.Println("mock: argocdapps uninstall finished")
}

var OpenStreamPlugin Plugin

func main() {
	fmt.Println("This is a plugin for OpenStream. Use it with OpenStream.")
}
