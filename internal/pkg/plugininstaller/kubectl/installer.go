package kubectl

import (
	"fmt"
	"os"
	"time"

	"github.com/cenkalti/backoff"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/util"
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

// WaitDeployReadyWithDeployList will wait all deploy in deployList get ready
func WaitDeployReadyWithDeployList(namespace string, deployList []string) plugininstaller.BaseOperation {
	waitFunc := func(options plugininstaller.RawOptions) error {
		operation := func() error {
			if err := util.CheckAllDeployAndServiceReady(namespace, deployList); err != nil {
				return err
			}
			return nil
		}
		bkoff := backoff.NewExponentialBackOff()
		bkoff.MaxElapsedTime = 3 * time.Minute
		err := backoff.Retry(operation, bkoff)
		if err != nil {
			return err
		}
		return nil
	}
	return waitFunc
}
