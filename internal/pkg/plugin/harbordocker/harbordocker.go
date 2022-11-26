package harbordocker

import (
	_ "embed"
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/docker/dockersh"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

const (
	HarborConfigFileName        = "harbor.yml"
	HarborScriptInstallFileName = "install.sh"
	HarborScriptCommonFileName  = "common.sh"
	HarborScriptPrepareFileName = "prepare"
)

//go:embed sh/harbor.tmpl.yml
var HarborConfigTemplate string

//go:embed sh/install.sh
var ScriptInstall string

//go:embed sh/common.sh
var ScriptCommon string

//go:embed sh/prepare
var ScriptPrepare string

var scripts = map[string]string{
	HarborScriptInstallFileName: ScriptInstall,
	HarborScriptCommonFileName:  ScriptCommon,
	HarborScriptPrepareFileName: ScriptPrepare,
}

func Install(options configmanager.RawOptions) error {
	if err := writeScripts(); err != nil {
		return err
	}

	// TODO(daniel-hutao): refactor is needed
	err := dockersh.ExecInSystemWithParams(".", []string{"./" + HarborScriptInstallFileName}, nil, true)
	if err != nil {
		return err
	}
	return nil
}

// renderConfig will render HarborConfigTemplate and then write it to disk.
func renderConfig(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts := Options{}
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	content, err := template.New().FromContent(HarborConfigTemplate).SetDefaultRender("compose", opts).Render()
	if err != nil {
		return nil, err
	}

	configFile, err := os.Create(HarborConfigFileName)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := configFile.Close(); err != nil {
			log.Errorf("Failed to close opened file (%s): %s.", configFile.Name(), err)
		}
	}()

	if _, err := configFile.Write([]byte(content)); err != nil {
		return nil, err
	}

	return options, err
}

func writeScripts() error {
	for name, sh := range scripts {
		err := os.WriteFile(name, []byte(sh), 0744)
		if err != nil {
			return err
		}
	}
	return nil
}
