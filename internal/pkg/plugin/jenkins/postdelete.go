package jenkins

import (
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func postDelete() error {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	if err = clearClusterRoleBinding(kubeClient); err != nil {
		return err
	}

	if err = clearClusterRole(kubeClient); err != nil {
		return err
	}

	if err = clearServiceAccount(kubeClient); err != nil {
		return err
	}

	clearPersistentVolume()

	return nil
}

func clearClusterRoleBinding(kubeClient *k8s.Client) error {
	return kubeClient.DeleteClusterRoleBinding(JenkinsName)
}

func clearClusterRole(kubeClient *k8s.Client) error {
	return kubeClient.DeleteClusterRole(JenkinsName)
}

func clearServiceAccount(kubeClient *k8s.Client) error {
	return kubeClient.DeleteServiceAccount(JenkinsName, JenkinsNamespace)
}

func clearPersistentVolume() {
	dataDir := getRealJenkinsDataDirectory()
	log.Warnf("\n\nNOTICE!!!\n"+
		"The PersistentVolume jenkins-pv is NOT been deleted.\n"+
		"You can execute the \"kubectl delete pv jenkins-pv\" to delete it manually."+
		"The local data directory %s is also NOT been deleted."+
		"You can delete it manually.\n", dataDir)
}
