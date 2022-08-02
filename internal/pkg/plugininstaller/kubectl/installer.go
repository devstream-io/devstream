package kubectl

import (
	"fmt"
	"os"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/kubectl"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// InstallByDownload will download file for apply
func ProcessByContent(action, downloadUrl, content string) plugininstaller.BaseOperation {
	return func(options plugininstaller.RawOptions) error {
		// generate k8s config file for apply
		configFileName, err := createKubectlFile(downloadUrl, content, options)
		if err != nil {
			return err
		}

		defer func() {
			err := os.Remove(configFileName)
			if err != nil {
				log.Debugf("kubectl delete temp file failed: %s", err)
			}
		}()
		// kubectl apply -f
		switch action {
		case "create":
			err = kubectl.KubeApply(configFileName)
		case "delete":
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
