package devlake

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/types"
)

const DevLakeSvcName = "devlake-lake"

// TODO(daniel-hutao): update the config below after devlake chart released.
var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "devlake/devlake",
		Timeout:     "5m",
		Wait:        types.Bool(true),
		UpgradeCRDs: types.Bool(true),
		ReleaseName: "devlake",
		Namespace:   "devlake",
	},
	Repo: helmCommon.Repo{
		URL:  "https://merico-dev.github.io/devlake-helm-chart",
		Name: "devlake",
	},
}

func genDevLakeState(options plugininstaller.RawOptions) (statemanager.ResourceState, error) {
	resState, err := helm.GetPluginAllState(options)
	if err != nil {
		return nil, err
	}

	// values.yaml
	opt, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}
	valuesYaml := opt.GetHelmParam().Chart.ValuesYaml
	resState["valuesYaml"] = valuesYaml

	// TODO(daniel-hutao): Use Ingress later.
	ip, err := getDevLakeClusterIP(opt.Chart.Namespace, DevLakeSvcName)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s:8080", ip)
	outputs := map[string]interface{}{
		"devlake_url": url,
	}
	resState.SetOutputs(outputs)

	return resState, nil
}

func getDevLakeClusterIP(namespace, name string) (string, error) {
	kClient, err := k8s.NewClient()
	if err != nil {
		return "", err
	}

	svc, err := kClient.GetService(namespace, name)
	if err != nil {
		return "", err
	}

	if svc.Spec.ClusterIP == "" {
		return "", fmt.Errorf("cluster ip is empty")
	}
	return svc.Spec.ClusterIP, nil
}
