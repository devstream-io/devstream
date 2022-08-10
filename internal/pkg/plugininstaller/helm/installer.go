package helm

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var (
	DefaultCreateOperations = plugininstaller.ExecuteOperations{
		DealWithNsWhenInstall,
		InstallOrUpdate,
	}
	DefaultUpdateOperations = plugininstaller.ExecuteOperations{
		InstallOrUpdate,
	}
	DefaultDeleteOperations = plugininstaller.ExecuteOperations{
		Delete,
		DealWithNsWhenInterruption,
	}
	DefaultTerminateOperations = plugininstaller.TerminateOperations{
		Delete,
		DealWithNsWhenInterruption,
	}
)

// InstallOrUpdate will install or update service by input options
func InstallOrUpdate(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	h, err := helm.NewHelm(opts.GetHelmParam())
	if err != nil {
		return err
	}

	log.Info("Creating or updating helm chart ...")
	if err := h.InstallOrUpgradeChart(); err != nil {
		log.Errorf("Failed to install or upgrade the chart: %s.", err)
		return err
	}
	return err
}

// DealWithNsWhenInstall will create namespace by input options
func DealWithNsWhenInstall(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	if !opts.CheckIfCreateNamespace() {
		log.Debugf("There's no need to delete the namespace for the create_namespace == false in the config file.")
		return nil
	}

	log.Debugf("Prepare to create the namespace: %s.", opts.GetNamespace())

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

// DealWithNsWhenInterruption will Delete namespace by input options
func DealWithNsWhenInterruption(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	if !opts.CreateNamespace {
		return err
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
	return err
}

// Delete will delete service base on input options
func Delete(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	h, err := helm.NewHelm(opts.GetHelmParam())
	if err != nil {
		return err
	}

	log.Infof("Uninstalling %s helm chart.", opts.GetReleaseName())
	if err = h.UninstallHelmChartRelease(); err != nil {
		return err
	}
	return nil
}
