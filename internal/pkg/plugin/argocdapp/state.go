package argocdapp

import (
	"time"

	"github.com/cenkalti/backoff"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/k8s"
)

func getStaticStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := newOptions(options)
	if err != nil {
		return nil, err
	}
	resStatus := make(statemanager.ResourceStatus)

	resStatus["app"] = map[string]interface{}{
		"name":      opts.App.Name,
		"namespace": opts.App.Namespace,
	}

	resStatus["src"] = map[string]interface{}{
		"repoURL":   opts.Source.RepoURL,
		"path":      opts.Source.Path,
		"valueFile": opts.Source.Valuefile,
	}

	resStatus["dest"] = map[string]interface{}{
		"server":    opts.Destination.Server,
		"namespace": opts.Destination.Namespace,
	}

	return resStatus, nil
}

func getDynamicStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := newOptions(options)
	if err != nil {
		return nil, err
	}

	retStatus := make(statemanager.ResourceStatus)
	operation := func() error {
		err := getArgoCDAppFromK8sAndSetStatus(retStatus, opts.App.Name, opts.App.Namespace)
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
	return retStatus, nil
}

func getArgoCDAppFromK8sAndSetStatus(status statemanager.ResourceStatus, name, namespace string) error {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	app, err := kubeClient.GetArgocdApplication(namespace, name)
	if err != nil {
		return err
	}

	d := kubeClient.DescribeArgocdApp(app)
	status["app"] = d["app"]
	status["src"] = d["src"]
	status["dest"] = d["dest"]

	return nil
}
