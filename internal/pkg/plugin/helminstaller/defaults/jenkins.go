package defaults

import (
	"fmt"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/helm"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	helmCommon "github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

var toolJenkins = "jenkins"

var DefaultConfigWithJenkins = helm.Options{
	Chart: helmCommon.Chart{
		ChartPath:   "",
		ChartName:   "jenkins/jenkins",
		Version:     "",
		Timeout:     "10m",
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

func init() {
	DefaultOptionsMap[toolJenkins] = &DefaultConfigWithJenkins
	StatusGetterFuncMap[toolJenkins] = GetJenkinsStatus
}

func GetJenkinsStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
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

	svcName, err := genJenkinsSvcName(options)
	if err != nil {
		return nil, err
	}

	// svc_name.svc_ns:svc_port
	url := fmt.Sprintf("http://%s.%s:8080", svcName, opt.Chart.Namespace)
	outputs := statemanager.ResourceOutputs{
		"jenkins_url": url,
	}

	resStatus.SetOutputs(outputs)

	return resStatus, nil
}

// See https://github.com/devstream-io/devstream/pull/1025#discussion_r952277174 and
// https://github.com/devstream-io/devstream/pull/1027#discussion_r953415932 for more info.
func genJenkinsSvcName(options configmanager.RawOptions) (string, error) {
	opts, err := helm.NewOptions(options)
	if err != nil {
		return "", err
	}

	pipe := func(s string) string {
		if len(s) > 63 {
			s = s[:63]
		}
		return strings.TrimSuffix(s, "-")
	}

	var tmpName string
	if strings.Contains(opts.Chart.ChartName, opts.Chart.ReleaseName) {
		tmpName = opts.Chart.ReleaseName
	} else {
		tmpName = fmt.Sprintf("%s-%s", opts.Chart.ReleaseName, opts.Chart.ChartName)
	}

	return pipe(tmpName), nil
}
