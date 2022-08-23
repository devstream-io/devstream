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

func genJenkinsState(options plugininstaller.RawOptions) (map[string]interface{}, error) {
	resource, err := helm.GetPluginAllState(options)
	if err != nil {
		return nil, err
	}

	// svc_name.svc_ns:svc_port
	resource["outputs"] = map[string]interface{}{
		"jenkins_url": "http://jenkins.jenkins:8080",
	}

	return resource, nil
}
