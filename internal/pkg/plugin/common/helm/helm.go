package helm

import (
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
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
type Options struct {
	CreateNamespace bool `mapstructure:"create_namespace"`
	Repo            helm.Repo
	Chart           helm.Chart
}

func (opts *Options) GetHelmParam() *helm.HelmParam {
	return &helm.HelmParam{
		Repo:  opts.Repo,
		Chart: opts.Chart,
	}
}

func InstallOrUpgradeChart(opts *Options) error {
	h, err := helm.NewHelm(opts.GetHelmParam())
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

func DealWithNsWhenInstall(opts *Options) error {
	if !opts.CreateNamespace {
		log.Debugf("There's no need to delete the namespace for the create_namespace == false in the config file.")
		return nil
	}

	log.Debugf("Prepare to create the namespace: %s.", opts.Chart.Namespace)

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	err = kubeClient.CreateNamespace(opts.Chart.Namespace)
	if err != nil {
		log.Debugf("Failed to create the namespace: %s.", opts.Chart.Namespace)
		return err
	}

	log.Debugf("The namespace %s has been created.", opts.Chart.Namespace)
	return nil
}

func DealWithNsWhenInterruption(opts *Options) error {
	if !opts.CreateNamespace {
		return nil
	}

	log.Debugf("Prepare to delete the namespace: %s.", opts.Chart.Namespace)

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	err = kubeClient.DeleteNamespace(opts.Chart.Namespace)
	if err != nil {
		log.Debugf("Failed to delete the namespace: %s.", opts.Chart.Namespace)
		return err
	}

	log.Debugf("The namespace %s has been deleted.", opts.Chart.Namespace)
	return nil
}
