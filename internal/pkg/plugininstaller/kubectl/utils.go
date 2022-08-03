package kubectl

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const defaultKubectlFileName = "kubectl-create-file_"

func writeContentToTmpFile(output *os.File, content string, opts *plugininstaller.RawOptions) error {
	t, err := template.New("app").Option("missingkey=error").Parse(content)
	if err != nil {
		return err
	}

	log.Debugf("All opts %+v", opts)

	err = t.Execute(output, opts)
	if err != nil {
		if strings.Contains(err.Error(), "can't evaluate field name") {
			msg := err.Error()
			start := strings.Index(msg, "<")
			end := strings.Index(msg, ">")
			return fmt.Errorf("plugin argocdapp needs options%s but it's missing from the config file", msg[start+1:end])
		} else {
			return fmt.Errorf("executing tpl error: %s", err)
		}
	}
	return nil
}

func createKubectlFile(downloadUrl, content string, options plugininstaller.RawOptions) (string, error) {
	// create temp file for kubectl apply
	tempFile, err := os.CreateTemp("", defaultKubectlFileName)
	if err != nil {
		return "", err
	}

	defer func() {
		err := tempFile.Close()
		if err != nil {
			log.Debugf("kubectl file close failed: %s", err)
		}
	}()

	if downloadUrl != "" {
		// download config file
		_, err := downloader.DownloadToFile(downloadUrl, tempFile)
		if err != nil {
			log.Debugf("Failed to download K8s deploy YAML file from %s.", downloadUrl)
			return "", err
		}
	} else if content != "" {
		// use content
		if err = writeContentToTmpFile(tempFile, content, &options); err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("No Install plugin config is set")
	}
	return tempFile.Name(), nil
}
