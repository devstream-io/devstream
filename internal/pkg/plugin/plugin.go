package plugin

import (
	"fmt"
	"os"
	"plugin"

	"github.com/merico-dev/stream/internal/pkg/config"
	"github.com/merico-dev/stream/internal/pkg/download"
)

// DevStreamPlugin is a struct, on which install/reinstall/uninstall interfaces are defined.
type DevStreamPlugin interface {
	// Install will return (true, nil) if there is no error occurred. Otherwise (false, error) will be returned.
	Install(*map[string]interface{}) (bool, error)
	Reinstall(*map[string]interface{}) (bool, error)
	Uninstall(*map[string]interface{}) (bool, error)
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

// Install loads the plugin and calls the Install method of that plugin.
func Install(tool *config.Tool) (bool, error) {
	path := fmt.Sprintf("plugins/%s_%s.so", tool.Name, tool.Version)
	appname := fmt.Sprintf("%s_%s.so", tool.Name, tool.Version)
	if !FileExist(path) {
		loader := download.NewDownloadClient()
		loader.AssetName = appname
		loader.Version = tool.Version
		loader.Filepath = path
		err := loader.GetAssetswithretry()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	p := loadPlugin(tool)
	return p.Install(&tool.Options)
}

// Reinstall loads the plugin and calls the Reinstall method of that plugin.
func Reinstall(tool *config.Tool) (bool, error) {
	p := loadPlugin(tool)
	return p.Reinstall(&tool.Options)
}

// Uninstall loads the plugin and calls the Uninstall method of that plugin.
func Uninstall(tool *config.Tool) (bool, error) {
	p := loadPlugin(tool)
	return p.Uninstall(&tool.Options)
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
