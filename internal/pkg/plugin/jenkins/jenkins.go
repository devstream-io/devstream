package jenkins

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"

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

// getHelmResourceAndCustomResource wraps helm resource and custom resource,
// this is due to the limitation of `plugininstaller`,
// now `plugininstaller.GetStateOperation` only support one resource get function,
// if we want to use both existing resource get function(such as helm's methods) and custom function,
// we have to wrap them into one function.
func getHelmResourceAndCustomResource(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	// 1. get helm resource
	resource, err := helm.GetPluginAllState(options)
	if err != nil {
		return nil, err
	}

	outputs := map[string]interface{}{}

	// TODO(daniel-hutao)
	outputs["foo"] = "bar"

	resource["outputs"] = outputs

	return resource, nil
}
