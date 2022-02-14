package devlake

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/downloader"
	"github.com/merico-dev/stream/pkg/util/kubectl"
)

func Uninstall(options *map[string]interface{}) (bool, error) {
	var param Param
	err := mapstructure.Decode(*options, &param)
	if err != nil {
		return false, err
	}

	if errs := validateParams(&param); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Param error: %s", e)
		}
		return false, fmt.Errorf("params are illegal")
	}

	_, err = downloader.Download(devLakeInstallYAMLDownloadURL, devLakeInstallYAMLFileName, ".")
	if err != nil {
		log.Debugf("Failed to download DevLake K8s deploy YAML file from %s", devLakeInstallYAMLDownloadURL)
		return false, err
	}

	err = kubectl.KubeDelete(devLakeInstallYAMLFileName)
	if err != nil {
		return false, err
	}
	if err = os.Remove(devLakeInstallYAMLFileName); err != nil {
		return false, err
	}

	return true, nil
}
