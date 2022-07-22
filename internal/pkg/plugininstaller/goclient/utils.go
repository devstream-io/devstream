package goclient

import (
	"strings"

	kerr "k8s.io/apimachinery/pkg/api/errors"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Check whether deployment is ready and service exists
func checkDeploymentsAndServicesReady(kubeClient *k8s.Client, opts *Options) (bool, error) {
	namespace := opts.Namespace
	deploy := opts.Deployment.Name
	svc := opts.Service.Name

	dp, err := kubeClient.GetDeployment(namespace, deploy)
	if err != nil {
		log.Debugf("Get deployment err: %s", err.Error())
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}

	if !kubeClient.IsDeploymentReady(dp) {
		log.Debugf("The deployment %s is not ready yet.", dp.Name)
		return false, nil
	}
	log.Debugf("The deployment %s is ready.", dp.Name)

	_, err = kubeClient.GetService(namespace, svc)
	if err != nil {
		log.Debugf("Get service err: %s", err.Error())
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	log.Debugf("The service %s is ready.", svc)

	return true, nil
}

// Delete application by goclient
func deleteApp(kubeClient *k8s.Client, opts *Options) error {
	// 1. Delete service
	if err := kubeClient.DeleteService(opts.Namespace, opts.Service.Name); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	// 2. Delete deployment
	if err := kubeClient.DeleteDeployment(opts.Namespace, opts.Deployment.Name); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}

	return nil
}

// Check whether the given namespace is created by dtm
// If the given namespace has label "created_by=DevStream", we'll control it.
// 1. The specified namespace is created by zentao plugin, then it should be deleted
//    when errors are encountered during creation or `dtm delete`.
// 2. The specified namespace is controlled by user, maybe they want to deploy zentao in
//    an existing namespace or other situations, then we should not delete this namespace.
func isDevstreamNSExists(kubeClient *k8s.Client, namespace string) (bool, error) {
	nsList, err := kubeClient.GetDevstreamNamespace()
	if err != nil {
		// not exist
		if kerr.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	for _, ns := range nsList.Items {
		// exist
		if ns.ObjectMeta.Name == namespace {
			return true, nil
		}
	}
	return false, nil
}
