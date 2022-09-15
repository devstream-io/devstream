package file

import (
	"io"
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
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		pluginDir = conf
	}

	if pluginDir == "" {
		return DefaultPluginDir, nil
	}

	pluginRealDir, err := getRealPath(pluginDir)
	if err != nil {
		return "", err
	}
	return pluginRealDir, nil
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

func CopyFile(srcFile, dstFile string) (err error) {
	// prepare source file
	sFile, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := sFile.Close(); closeErr != nil {
			log.Errorf("Failed to close file %s: %s", srcFile, closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// create destination file
	dFile, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dFile.Close(); closeErr != nil {
			log.Errorf("Failed to close file %s: %s", dstFile, closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// copy and sync
	if _, err = io.Copy(dFile, sFile); err != nil {
		return nil
	}
	return dFile.Sync()
}
