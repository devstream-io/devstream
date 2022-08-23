package jenkins

import (
	"fmt"
	"strings"

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

	// values.yaml
	opt, err := helm.NewOptions(options)
	if err != nil {
		return nil, err
	}
	valuesYaml := opt.GetHelmParam().Chart.ValuesYaml
	resState["values_yaml"] = valuesYaml

	svcName, err := genJenkinsSvcName(options)
	if err != nil {
		return nil, err
	}

	// svc_name.svc_ns:svc_port
	url := fmt.Sprintf("http://%s.%s:8080", svcName, opt.Chart.Namespace)
	outputs := map[string]interface{}{
		"jenkins_url": url,
	}

	resState.SetOutputs(outputs)

	return resState, nil
}

// see https://github.com/devstream-io/devstream/pull/1025#discussion_r952277174 for more info.
func genJenkinsSvcName(options plugininstaller.RawOptions) (string, error) {
	opts, err := helm.NewOptions(options)
	if err != nil {
		return "", err
	}

	var retStr string
	if strings.Contains(opts.Chart.ChartName, opts.Chart.ReleaseName) {
		retStr = strings.TrimSuffix(opts.Chart.ReleaseName, "-")
		if len(retStr) > 63 {
			retStr = retStr[:64]
		}
		return retStr, nil
	}
	retStr = fmt.Sprintf("%s-%s", opts.Chart.ReleaseName, opts.Chart.ChartName)

	return retStr, nil
}
