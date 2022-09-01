package file

import (
	"os"
	"path/filepath"
	"strings"

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

func GetPluginDir(conf string) (string, error) {
	if flag := viper.GetString("plugin-dir"); flag != "" {
		return flag, nil
	}

	if conf == "" {
		return DefaultPluginDir, nil
	}

	pluginDir, err := getRealPath(conf)
	if err != nil {
		return "", err
	}
	return pluginDir, nil
}

// getRealPath deal with "~" in the filePath
func getRealPath(filePath string) (string, error) {
	if !strings.Contains(filePath, "~") {
		return filePath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	retPath := filepath.Join(homeDir, strings.TrimPrefix(filePath, "~"))
	log.Debugf("real path: %s.", retPath)

	return retPath, nil
}

func SetPluginDir(conf string) error {
	pluginDir, err := GetPluginDir(conf)
	if err != nil {
		return err
	}
	viper.Set("plugin-dir", pluginDir)
	return nil
}

// CopyFile will copy file content from src to dst
func CopyFile(src, dest string) error {
	bytesRead, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, bytesRead, 0644)
}
