package kubectl

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/kubectl"
)

// InstallByDownload will download file for apply
func ProcessByContent(action string, templateConfig *file.TemplateConfig) plugininstaller.BaseOperation {
	return func(options plugininstaller.RawOptions) error {
		// generate k8s config file for apply
		configFileName, err := templateConfig.RenderFile("kubectl", options).Run()
		if err != nil {
			return err
		}
		// kubectl apply -f
		switch action {
		case kubectl.APPLY:
			err = kubectl.KubeApply(configFileName)
		case kubectl.CREATE:
			err = kubectl.KubeCreate(configFileName)
		case kubectl.DELETE:
			err = kubectl.KubeDelete(configFileName)
		default:
			err = fmt.Errorf("kubectl not support this kind of action: %s", action)
		}
		if err != nil {
			return err
		}
		return nil
	}
}
