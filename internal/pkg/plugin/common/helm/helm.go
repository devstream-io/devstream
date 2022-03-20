package helm

import (
	"github.com/merico-dev/stream/pkg/util/helm"
	"github.com/merico-dev/stream/pkg/util/k8s"
	"github.com/merico-dev/stream/pkg/util/log"
)

// NOTICE: Don't use:
// type Param struct {
// 	CreateNamespace bool `mapstructure:"create_namespace"`
// 	helm.HelmParam
// }
// or
// type Param struct {
// 	CreateNamespace bool `mapstructure:"create_namespace"`
// 	*helm.HelmParam
// }
// see pr #174 for more info

// Param is the struct for parameters used by the argocd package.
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

func InstallOrUpgradeChart(param *Param) error {
	h, err := helm.NewHelm(param.GetHelmParam())
	if err != nil {
		return err
	}

	log.Info("Creating or updating helm chart ...")
	if err := h.InstallOrUpgradeChart(); err != nil {
		log.Debugf("Failed to install or upgrade the chart: %s.", err)
		return err
	}
	return nil
}

func DealWithNsWhenInstall(param *Param) error {
	if !param.CreateNamespace {
		log.Debugf("There's no need to delete the namespace for the create_namespace == false in the config file.")
		return nil
	}

	log.Debugf("Prepare to create the namespace: %s.", param.Chart.Namespace)

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	err = kubeClient.CreateNamespace(param.Chart.Namespace)
	if err != nil {
		log.Debugf("Failed to create the namespace: %s.", param.Chart.Namespace)
		return err
	}

	log.Debugf("The namespace %s has been created.", param.Chart.Namespace)
	return nil
}

func DealWithNsWhenInterruption(param *Param) error {
	if !param.CreateNamespace {
		return nil
	}

	log.Debugf("Prepare to delete the namespace: %s.", param.Chart.Namespace)

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	err = kubeClient.DeleteNamespace(param.Chart.Namespace)
	if err != nil {
		log.Debugf("Failed to delete the namespace: %s.", param.Chart.Namespace)
		return err
	}

	log.Debugf("The namespace %s has been deleted.", param.Chart.Namespace)
	return nil
}
