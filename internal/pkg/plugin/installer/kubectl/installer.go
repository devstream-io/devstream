package kubectl

import (
	"fmt"
	"io"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/pkg/util/kubectl"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/pkgerror"
	"github.com/devstream-io/devstream/pkg/util/template"
)

func ProcessByContent(action, content string) installer.BaseOperation {
	return func(options configmanager.RawOptions) error {
		reader, err := renderKubectlContent(content, options)
		if err != nil {
			return err
		}

		return processByIOReader(action, reader)
	}
}

func renderKubectlContent(content string, options configmanager.RawOptions) (io.Reader, error) {
	content, err := template.New().FromContent(content).SetDefaultRender("kubectl", options).Render()
	if err != nil {
		return nil, err
	}
	if content == "" {
		return nil, fmt.Errorf("kubectl content is empty")
	}

	return strings.NewReader(content), nil
}

func ProcessByURL(action, url string) installer.BaseOperation {
	return func(options configmanager.RawOptions) error {
		content, err := template.New().FromURL(url).SetDefaultRender("kubectl", options).Render()
		if err != nil {
			return err
		}
		if content == "" {
			return fmt.Errorf("kubectl content is empty")
		}

		reader := strings.NewReader(content)

		return processByIOReader(action, reader)
	}
}

func processByIOReader(action string, reader io.Reader) error {
	// generate k8s config file for apply
	var err error
	// kubectl apply -f
	switch action {
	case kubectl.Create:
		err = kubectl.KubeApplyFromIOReader(reader)
	case kubectl.Apply:
		err = kubectl.KubeApplyFromIOReader(reader)
	case kubectl.Delete:
		err = kubectl.KubeDeleteFromIOReader(reader)
		// ignore resource not exist error
		if err != nil && pkgerror.CheckErrorMatchByMessage(err, kubectl.ArgocdApplicationNotExist) {
			log.Warnf("kubectl config is already deleted")
			err = nil
		}
	default:
		err = fmt.Errorf("kubectl not support this kind of action: %v", action)
	}
	if err != nil {
		return err
	}
	return nil
}
