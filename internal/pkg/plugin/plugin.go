package plugin

import (
	"fmt"
	"os"
	"plugin"

	"github.com/ironcore864/openstream/internal/pkg/config"
)

type OpenStreamPlugin interface {
	Install()
	Reinstall()
	Uninstall()
}

func loadPlugin(tool *config.Tool) OpenStreamPlugin {
	mod := fmt.Sprintf("plugins/%s_%s.so", tool.Name, tool.Version)
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var openStreamPlugin OpenStreamPlugin
	symOpenStreamPlugin, err := plug.Lookup("OpenStreamPlugin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	openStreamPlugin, ok := symOpenStreamPlugin.(OpenStreamPlugin)
	if !ok {
		fmt.Println(err)
		os.Exit(1)
	}

	return openStreamPlugin
}

func Install(tool *config.Tool) {
	plugin := loadPlugin(tool)
	plugin.Install()
}

func Reinstall(tool *config.Tool) {
	plugin := loadPlugin(tool)
	plugin.Reinstall()
}

func Uninstall(tool *config.Tool) {
	plugin := loadPlugin(tool)
	plugin.Uninstall()
}
