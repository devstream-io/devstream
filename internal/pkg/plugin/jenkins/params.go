package jenkins

import "github.com/merico-dev/stream/pkg/util/helm"

// Param is the struct for parameters used by the jenkins package.
type Param struct {
	CreateNamespace bool `mapstructure:"create_namespace"`
	Repo            helm.Repo
	Chart           helm.Chart
}

func (p *Param) GetHelmParam() *helm.HelmParam {
	return &helm.HelmParam{
		Repo:  p.Repo,
		Chart: p.Chart,
	}
}

func (p *Param) renderValuesYamlForJenkins() {
	p.Chart.ValuesYaml = `persistence:
  storageClass: jenkins-pv
serviceAccount:
  create: false
  name: jenkins
`
}
