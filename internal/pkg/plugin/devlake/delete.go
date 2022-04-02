package devlake

import (
	"os"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/kubectl"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Delete(options map[string]interface{}) (bool, error) {
	var opts Options

	// decode input parameters into a struct
	err := mapstructure.Decode(options, &opts)
	if err != nil {
		return false, err
	}

	// download DevLake installation YAML file
	_, err = downloader.Download(devLakeInstallYAMLDownloadURL, devLakeInstallYAMLFileName, ".")
	if err != nil {
		log.Debugf("Failed to download DevLake K8s deploy YAML file from %s.", devLakeInstallYAMLDownloadURL)
		return false, err
	}

	// kubectl delete -f
	err = kubectl.KubeDelete(devLakeInstallYAMLFileName)
	if err != nil {
		return false, err
	}

	// remove temporary YAML file used for kubectl apply
	if err = os.Remove(devLakeInstallYAMLFileName); err != nil {
		log.Warnf("Temporary YAML file %s can't be deleted, but the installation is successful.", devLakeInstallYAMLFileName)
	}

	return true, nil
}
