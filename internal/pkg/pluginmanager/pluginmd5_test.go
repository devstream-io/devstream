package pluginmanager

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/merico-dev/stream/cmd/devstream/version"
	"github.com/merico-dev/stream/internal/pkg/configloader"
)

// TestCheckLocalPlugins test plguin .so matches dtm core md5
func TestCheckLocalPlugins(t *testing.T) {
	viper.Set("plugin-dir", "./")

	tools := []configloader.Tool{
		{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
	}

	file := configloader.GetPluginFileName(&tools[0])
	fileMD5 := configloader.GetPluginMD5FileName(&tools[0])

	config := &configloader.Config{Tools: tools}

	err := createNewFile(file)
	assert.NoError(t, err)

	err = addMD5File(file, fileMD5)
	assert.NoError(t, err)

	err = CheckLocalPlugins(config)
	assert.NoError(t, err)

}

// TestCheckPluginMismatch test checkPluginMismatch error
func TestCheckPluginMismatch(t *testing.T) {
	viper.Set("plugin-dir", "./")

	tools := []configloader.Tool{
		{Name: "c", Plugin: configloader.Plugin{Kind: "c"}},
	}
	file := configloader.GetPluginFileName(&tools[0])
	fileMD5 := configloader.GetPluginMD5FileName(&tools[0])

	err := createNewFile(file)
	assert.NoError(t, err)

	err = createNewFile(fileMD5)
	assert.NoError(t, err)

	err = checkPluginMismatch(viper.GetString("plugin-dir"), file, fileMD5, tools[0].Name)
	expectErrMsg := fmt.Sprintf("plugin %s doesn't match with .md5", tools[0].Name)
	assert.EqualError(t, err, expectErrMsg)
}

func createNewFile(fileName string) error {
	f1, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f1.Close()
	return nil
}

func addMD5File(fileName, md5FileName string) error {
	md5, err := version.CalcFileMD5(fileName)
	if err != nil {
		return err
	}
	md5File, err := os.Create(md5FileName)
	if err != nil {
		return err
	}
	_, err = md5File.Write([]byte(md5))
	if err != nil {
		return err
	}
	return nil
}
