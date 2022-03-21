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
	clearMD5string()
	viper.Set("plugin-dir", "./")

	tools := []configloader.Tool{
		{Name: "a", Plugin: configloader.Plugin{Kind: "a"}},
		{Name: "b", Plugin: configloader.Plugin{Kind: "b"}},
	}

	fileA := configloader.GetPluginFileName(&tools[0])
	fileB := configloader.GetPluginFileName(&tools[1])

	defer func() { clearMD5string() }()

	config := &configloader.Config{Tools: tools}

	err := createNewFile(fileA)
	assert.NoError(t, err)

	err = addMD5string(fileA)
	assert.NoError(t, err)

	err = createNewFile(fileB)
	assert.NoError(t, err)

	err = addMD5string(fileB)
	assert.NoError(t, err)

	err = CheckLocalPlugins(config)
	assert.NoError(t, err)

}

// TestCheckPluginMismatch test checkPluginMismatch error
func TestCheckPluginMismatch(t *testing.T) {
	clearMD5string()
	viper.Set("plugin-dir", "./")

	tools := []configloader.Tool{
		{Name: "c", Plugin: configloader.Plugin{Kind: "c"}},
	}
	defer func() { clearMD5string() }()
	fileC := configloader.GetPluginFileName(&tools[0])

	err := createNewFile(fileC)
	assert.NoError(t, err)

	err = checkPluginMismatch(viper.GetString("plugin-dir"), fileC, tools[0].Name)
	expectErrMsg := fmt.Sprintf("plugin %s doesn't match with dtm core", tools[0].Name)
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

func addMD5string(fileName string) error {
	md5, err := version.CalcFileMD5(fileName)
	if err != nil {
		return err
	}
	if version.MD5String == "" {
		version.MD5String = md5
	} else {
		version.MD5String = fmt.Sprintf("%s:%s", version.MD5String, md5)
	}
	return nil
}

func clearMD5string() {
	version.MD5String = ""
}
