package plugin

import (
	"fmt"
	"os"
	"plugin"
)

type OpenStreamPlugin interface {
	Install()
	Reinstall()
	Uninstall()
}

func loadPlugin(pkg string) OpenStreamPlugin {
	pluginPathName := fmt.Sprintf("./plugins/%s.so", pkg)
	plug, err := plugin.Open(pluginPathName)
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

func Install(pkg string) {
	plugin := loadPlugin(pkg)
	plugin.Install()
}

func Reinstall(pkg string) {
	plugin := loadPlugin(pkg)
	plugin.Reinstall()
}

func Uninstall(pkg string) {
	plugin := loadPlugin(pkg)
	plugin.Uninstall()
}
