package jenkins

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var defaultHelmConfig = helm.Options{
	Chart: helmCommon.Chart{
		ChartName:   "jenkins/jenkins",
		Timeout:     "5m",
		UpgradeCRDs: types.Bool(true),
		Wait:        types.Bool(true),
		ReleaseName: "jenkins",
		Namespace:   "jenkins",
	},
	Repo: helmCommon.Repo{
		URL:  "https://charts.jenkins.io",
		Name: "jenkins",
	},
}

func genJenkinsState(options plugininstaller.RawOptions) (statemanager.ResourceState, error) {
	resState, err := helm.GetPluginAllState(options)
	if err != nil {
		return nil, err
	}

	// svc_name.svc_ns:svc_port
	outputs := map[string]interface{}{
		"jenkins_url": "http://jenkins.jenkins:8080",
	}
	resState.SetOutputs(outputs)

	// values.yaml
	opt, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}

	valuesYaml := opt.GetHelmParam().Chart.ValuesYaml
	resState["values_yaml"] = valuesYaml

	return resState, nil
}
