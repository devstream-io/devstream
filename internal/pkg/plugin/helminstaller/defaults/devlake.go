package defaults

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/types"
)

const (
	DevLakeSvcName = "devlake-lake"
	toolDevLake    = "devlake"
)

var DefaultConfigWithDevLake = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "devlake/devlake",
		Version:     "",
		Timeout:     "10m",
		Wait:        types.Bool(true),
		UpgradeCRDs: types.Bool(true),
		ReleaseName: "devlake",
		Namespace:   "devlake",
	},
	Repo: helmCommon.Repo{
		URL:  "https://apache.github.io/incubator-devlake-helm-chart",
		Name: "devlake",
	},
}

func init() {
	RegisterDefaultHelmAppInstance(toolDevLake, &DefaultConfigWithDevLake, GetDevLakeStatus)
}

func GetDevLakeStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	resStatus, err := helm.GetAllResourcesStatus(options)
	if err != nil {
		return nil, err
	}

	// values.yaml
	opt, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}
	valuesYaml := opt.GetHelmParam().Chart.ValuesYaml
	resStatus["valuesYaml"] = valuesYaml

	// TODO(daniel-hutao): Use Ingress later.
	ip, err := getDevLakeClusterIP(opt.Chart.Namespace, DevLakeSvcName)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s:8080", ip)
	outputs := statemanager.ResourceOutputs{
		"devlake_url": url,
	}
	resStatus.SetOutputs(outputs)

	return resStatus, nil
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
