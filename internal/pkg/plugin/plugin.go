package plugin

import (
	"fmt"
	"os"
	"plugin"

	"github.com/merico-dev/stream/internal/pkg/config"
)

type DevStreamPlugin interface {
	Install(*map[string]interface{})
	Reinstall(*map[string]interface{})
	Uninstall(*map[string]interface{})
}

func loadPlugin(tool *config.Tool) DevStreamPlugin {
	mod := fmt.Sprintf("plugins/%s_%s.so", tool.Name, tool.Version)
	plug, err := plugin.Open(mod)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var devStreamPlugin DevStreamPlugin
	symDevStreamPlugin, err := plug.Lookup("DevStreamPlugin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	devStreamPlugin, ok := symDevStreamPlugin.(DevStreamPlugin)
	if !ok {
		fmt.Println(err)
		os.Exit(1)
	}

	return devStreamPlugin
}

func Install(tool *config.Tool) {
	p := loadPlugin(tool)
	p.Install(&tool.Options)
}

func Reinstall(tool *config.Tool) {
	p := loadPlugin(tool)
	p.Reinstall(&tool.Options)
}

func Uninstall(tool *config.Tool) {
	p := loadPlugin(tool)
	p.Uninstall(&tool.Options)
}
