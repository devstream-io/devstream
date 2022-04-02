package devlake

import (
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/kubectl"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options

	// decode input parameters into a struct
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return nil, err
	}

	// download DevLake installation YAML file
	_, err = downloader.Download(devLakeInstallYAMLDownloadURL, devLakeInstallYAMLFileName, ".")
	if err != nil {
		log.Debugf("Failed to download DevLake K8s deploy YAML file from %s.", devLakeInstallYAMLDownloadURL)
		return nil, err
	}

	// kubectl apply -f
	if err = kubectl.KubeApply(devLakeInstallYAMLDownloadURL); err != nil {
		return nil, err
	}

	// remove temporary YAML file used for kubectl apply
	if err = os.Remove(devLakeInstallYAMLFileName); err != nil {
		log.Warnf("Temporary YAML file %s can't be deleted, but the installation is successful.", devLakeInstallYAMLFileName)
	}

	// wait till deployments are ready
	operation := func() error {
		if err := allDeploymentsAndServicesReady(); err != nil {
			return err
		}
		return nil
	}
	bkoff := backoff.NewExponentialBackOff()
	bkoff.MaxElapsedTime = 3 * time.Minute
	err = backoff.Retry(operation, bkoff)
	if err != nil {
		return nil, err
	}

	// build state & return results
	return buildState(opts), nil
}
