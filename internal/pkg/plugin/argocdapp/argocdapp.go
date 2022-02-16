package argocdapp

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/k8s"
)

const argoCDAppYAMLFile = "./app.yaml"

func writeContentToTmpFile(file string, content string, param *Param) error {
	t, err := template.New("app").Option("missingkey=error").Parse(content)
	if err != nil {
		return err
	}

	output, err := os.Create(file)
	if err != nil {
		return err
	}

	err = t.Execute(output, param)
	if err != nil {
		if strings.Contains(err.Error(), "can't evaluate field name") {
			msg := err.Error()
			start := strings.Index(msg, "<")
			end := strings.Index(msg, ">")
			return fmt.Errorf("plugin argocdapp needs options%s but it's missing from the config file", msg[start+1:end])
		} else {
			return fmt.Errorf("executing tpl: %s", err)
		}
	}
	return nil
}

func buildState(p Param) map[string]interface{} {
	res := make(map[string]interface{})

	res["app"] = map[string]interface{}{
		"name":      p.App.Name,
		"namespace": p.App.Namespace,
	}

	res["src"] = map[string]interface{}{
		"repoURL":   p.Source.RepoURL,
		"path":      p.Source.Path,
		"valueFile": p.Source.Valuefile,
	}

	res["dest"] = map[string]interface{}{
		"server":    p.Destination.Server,
		"namespace": p.Destination.Namespace,
	}

	return res
}

// isArgoCDAppReady returns nil if the app is ready; otherwise it returns an error
func isArgoCDAppReady(name, namespace string) error {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	app, err := kubeClient.GetArgocdApplication(namespace, name)
	if err != nil {
		return err
	}

	if kubeClient.IsArgocdApplicationReady(app) {
		log.Infof("%s/%s is ready", namespace, name)
		return nil
	} else {
		log.Infof("%s/%s is not ready yet", namespace, name)
		return fmt.Errorf("%s/%s not ready", namespace, name)
	}
}

func getArgoCDAppFromK8sAndSetState(state map[string]interface{}, name, namespace string) error {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	app, err := kubeClient.GetArgocdApplication(namespace, name)
	if err != nil {
		return err
	}

	state = kubeClient.DescribeArgocdApp(app)
	return nil
}
