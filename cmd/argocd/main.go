package main

import (
	"fmt"
	"log"

	"github.com/ironcore864/openstream/internal/pkg/argocd"
)

const NAME = "argocd"

type Plugin string

func (p Plugin) Install(options *map[string]interface{}) {
	argocd.Install(options)
	log.Println("argocd install finished")
}

func (p Plugin) Reinstall(options *map[string]interface{}) {
	log.Println("mock: argocd reinstall finished")
}

func (p Plugin) Uninstall(options *map[string]interface{}) {
	log.Println("mock: argocd uninstall finished")
}

var OpenStreamPlugin Plugin

func main() {
	fmt.Println("This is a plugin for OpenStream. Use it with OpenStream.")
}
