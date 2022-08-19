package file

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/pkg/util/log"
)

var DefaultPluginDir string

func init() {
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("failed to get home dir: %s", err)
	}
	DefaultPluginDir = filepath.Join(homeDir, ".devstream", "plugins")
}

func GetPluginDir(conf string) string {
	if flag := viper.GetString("plugin-dir"); flag != "" {
		return flag
	}
	if conf != "" {
		return conf
	}
	return DefaultPluginDir
}

func SetPluginDir(conf string) {
	viper.Set("plugin-dir", GetPluginDir(conf))
}

// CopyFile will copy file content from src to dst
func CopyFile(src, dest string) error {
	bytesRead, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, bytesRead, 0644)
}
