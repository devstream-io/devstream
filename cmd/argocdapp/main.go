package main

import (
	"fmt"
	"log"

	"github.com/merico-dev/stream/internal/pkg/argocdapp"
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

var DevStreamPlugin Plugin

func main() {
	fmt.Println("This is a plugin for DevStream. Use it with DevStream.")
}
