package argocdapp

import (
	"time"

	"github.com/cenkalti/backoff"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/k8s"
)

func getStaticState(options plugininstaller.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})

	res["app"] = map[string]interface{}{
		"name":      opts.App.Name,
		"namespace": opts.App.Namespace,
	}

	res["src"] = map[string]interface{}{
		"repoURL":   opts.Source.RepoURL,
		"path":      opts.Source.Path,
		"valueFile": opts.Source.Valuefile,
	}

	res["dest"] = map[string]interface{}{
		"server":    opts.Destination.Server,
		"namespace": opts.Destination.Namespace,
	}

	return res, nil
}

func getDynamicState(options plugininstaller.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	state := make(map[string]interface{})
	operation := func() error {
		err := getArgoCDAppFromK8sAndSetState(state, opts.App.Name, opts.App.Namespace)
		if err != nil {
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
	return state, nil
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

	d := kubeClient.DescribeArgocdApp(app)
	state["app"] = d["app"]
	state["src"] = d["src"]
	state["dest"] = d["dest"]

	return nil
}
