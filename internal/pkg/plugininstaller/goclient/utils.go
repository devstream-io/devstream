package goclient

import (
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var checkRetryTime = 5

// Check whether deployment is ready and service exists
func checkDeploymentsAndServicesReady(kubeClient k8s.K8sAPI, opts *Options) (bool, error) {
	namespace := opts.Namespace
	deploy := opts.Deployment.Name
	svc := opts.Service.Name

	dp, err := kubeClient.GetDeployment(namespace, deploy)
	if err != nil {
		log.Debugf("Get deployment err: %s", err.Error())
		if kerr.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	err = kubeClient.WaitForDeploymentReady(checkRetryTime, opts.Namespace, opts.Deployment.Name)
	if err != nil {
		log.Debugf("The deployment %s is not ready yet.", dp.Name)
		return false, nil
	}
	log.Debugf("The deployment %s is ready.", dp.Name)

	_, err = kubeClient.GetService(namespace, svc)
	if err != nil {
		log.Debugf("Get service err: %s", err.Error())
		if kerr.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	log.Debugf("The service %s is ready.", svc)

	return true, nil
}

// Delete application by goclient
func deleteApp(kubeClient k8s.K8sAPI, opts *Options) error {
	// 1. Delete service
	if err := kubeClient.DeleteService(opts.Namespace, opts.Service.Name); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}

	// 2. Delete deployment
	if err := kubeClient.DeleteDeployment(opts.Namespace, opts.Deployment.Name); err != nil {
		if !kerr.IsNotFound(err) {
			return err
		}
	}

	return nil
}

// Generate []corev1.Volume for deployment from Options.PersistentVolumeClaims
func (opts *Options) genVolumesForDeployment() []corev1.Volume {
	var v []corev1.Volume
	for _, pvc := range opts.PersistentVolumeClaims {
		v = append(v, corev1.Volume{
			Name: pvc.PVCName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvc.PVCName,
				},
			},
		})
	}

	return v
}

// Generate []corev1.EnvVar for deployment from Options.PersistentVolumeClaims
func (opts *Options) genEnvsForDeployment() []corev1.EnvVar {
	var e []corev1.EnvVar
	for _, env := range opts.Deployment.Envs {
		e = append(e, corev1.EnvVar{
			Name:  env.Key,
			Value: env.Value,
		})
	}
	return e
}
