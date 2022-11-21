package goclient

import (
	"errors"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// Create namespace by goclient
func DealWithNsWhenInstall(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// Namespace should not be empty
	if opts.Namespace == "" {
		log.Debugf("No namespace is given.")
		return errors.New("No namespace is given.")
	}
	log.Debugf("Prepare to create the namespace: %s.", opts.Namespace)

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	// Check whether the given namespace already exists.
	exist, err := kubeClient.IsNamespaceExists(opts.Namespace)
	if err != nil {
		log.Debugf("Failed to check whether namespace: %s exists.", opts.Namespace)
		return err
	}
	if exist {
		log.Debugf("Namespace: %s already exists.", opts.Namespace)
		return nil
	}

	// Create new namespace
	if err = kubeClient.CreateNamespace(opts.Namespace); err != nil {
		log.Debugf("Failed to create the namespace: %s.", opts.Namespace)
		return err
	}

	log.Debugf("The namespace %s has been created.", opts.Namespace)
	return nil
}

// Deal with resource when errors occur during creation
func DealWithErrWhenInstall(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	exist, err := kubeClient.IsDevstreamNS(opts.Namespace)
	if err != nil {
		log.Debugf("Failed to check whether namespace: %s exists.", opts.Namespace)
		return err
	}

	// namespace is controlled by dtm, just delete this namespace
	if exist {
		log.Debugf("Prepare to delete the namespace: %s.", opts.Namespace)

		err := kubeClient.DeleteNamespace(opts.Namespace)
		if err != nil {
			log.Debugf("Failed to delete the namespace: %s.", opts.Namespace)
			return err
		}

		log.Debugf("The namespace %s has been deleted.", opts.Namespace)
	}

	return nil
}

// Create persistent volume with hostpath
func CreatePersistentVolumeWrapper(pvPath map[string]string) installer.BaseOperation {
	return func(options configmanager.RawOptions) error {
		opts, err := NewOptions(options)
		if err != nil {
			return err
		}

		kubeClient, err := k8s.NewClient()
		if err != nil {
			return err
		}

		for _, pv := range opts.PersistentVolumes {
			newPVOpt := &k8s.PVOption{
				Name:             pv.PVName,
				StorageClassName: opts.StorageClassName,
				AccessMode: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Capacity:                      pv.PVCapacity,
				PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
			}

			if path, ok := pvPath[newPVOpt.Name]; ok {
				newPVOpt.HostPath = path
			}

			log.Debugf("Prepare to create persistentVolume: %s.", newPVOpt.Name)

			if err := kubeClient.CreatePersistentVolume(newPVOpt); err != nil {
				if !kerr.IsAlreadyExists(err) {
					return err
				}
				log.Infof("The PersistentVolume %s is already exists.", newPVOpt.Name)
			}
			log.Debugf("The PersistentVolume %s has been created.", newPVOpt.Name)
		}

		return nil
	}
}

// Create persistent volume claim
func CreatePersistentVolumeClaim(options configmanager.RawOptions) error {

	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	for _, pvc := range opts.PersistentVolumeClaims {
		newPVOpt := &k8s.PVCOption{
			Name:             pvc.PVCName,
			NameSpace:        opts.Namespace,
			StorageClassName: opts.StorageClassName,
			AccessMode: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Requirement: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(pvc.PVCCapacity),
				},
			},
		}

		log.Debugf("Prepare to create persistentVolumeClaim: %s.", newPVOpt.Name)
		if err := kubeClient.CreatePersistentVolumeClaim(newPVOpt); err != nil {
			if !kerr.IsAlreadyExists(err) {
				return err
			}
			log.Infof("The PersistentVolumeClaim %s is already exists.", newPVOpt.Name)
		}
		log.Debugf("The PersistentVolumeClaim %s has been created.", newPVOpt.Name)

	}

	return nil
}

// Create deployment by goclient with label, containerPorts and name
func CreateDeploymentWrapperLabelAndContainerPorts(label map[string]string, containerPorts *[]corev1.ContainerPort, name string) installer.BaseOperation {
	return func(options configmanager.RawOptions) error {

		opts, err := NewOptions(options)
		if err != nil {
			return err
		}

		kubeClient, err := k8s.NewClient()
		if err != nil {
			return err
		}

		volumes := opts.genVolumesForDeployment()
		envs := opts.genEnvsForDeployment()

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:   opts.Deployment.Name,
				Labels: label,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: pointer.Int32Ptr(int32(opts.Deployment.Replicas)),
				Selector: &metav1.LabelSelector{MatchLabels: label},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: label,
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  name,
								Image: opts.Deployment.Image,
								Env:   envs,
								Ports: *containerPorts,
							},
						},
						Volumes: volumes,
					},
				},
			},
		}

		log.Debugf("The Deployment %s has been created.", deployment.Name)
		if err := kubeClient.CreateDeployment(opts.Namespace, deployment); err != nil {
			if !kerr.IsAlreadyExists(err) {
				return err
			}
			log.Infof("The Deployment %s is already exists.", deployment.Name)
		}
		log.Debugf("The Deployment %s has been created.", deployment.Name)

		return nil
	}
}

// Create service by goclient with label and servicePort
func CreateServiceWrapperLabelAndPorts(label map[string]string, svcPort *corev1.ServicePort) installer.BaseOperation {
	return func(options configmanager.RawOptions) error {

		opts, err := NewOptions(options)
		if err != nil {
			return err
		}

		kubeClient, err := k8s.NewClient()
		if err != nil {
			return err
		}

		svcPort.NodePort = int32(opts.Service.NodePort)
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:   opts.Service.Name,
				Labels: label,
			},
			Spec: corev1.ServiceSpec{
				Ports:    []corev1.ServicePort{*svcPort},
				Selector: label,
				Type:     corev1.ServiceTypeNodePort,
			},
		}

		log.Debugf("The Service %s has been created.", svc.Name)
		return kubeClient.CreateService(opts.Namespace, svc)
	}
}

// Check application status by goclient with retry times
func WaitForReady(retry int) installer.BaseOperation {
	return func(options configmanager.RawOptions) error {

		opts, err := NewOptions(options)
		if err != nil {
			return err
		}

		kubeClient, err := k8s.NewClient()
		if err != nil {
			return err
		}

		err = kubeClient.WaitForDeploymentReady(retry, opts.Namespace, opts.Deployment.Name)
		if err != nil {
			return err
		}
		return nil
	}
}

// Delete plugin by goclient
func Delete(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. Create k8s clientset
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	// 2. Delete application
	if err = deleteApp(kubeClient, &opts); err != nil {
		return err
	}

	// 3. Delete PVC
	log.Debug("Prepare to delete PVC.")
	for _, pvc := range opts.PersistentVolumeClaims {
		if err = kubeClient.DeletePersistentVolumeClaim(opts.Namespace, pvc.PVCName); err != nil {
			if !strings.Contains(err.Error(), "not found") {
				return err
			}
		}
	}

	// 4. Delete PV
	log.Debug("Prepare to delete PV.")
	for _, pv := range opts.PersistentVolumes {
		if err = kubeClient.DeletePersistentVolume(pv.PVName); err != nil {
			if !strings.Contains(err.Error(), "not found") {
				return err
			}
		}
	}

	// 5. Delete namespace only when namespace is controlled by dtm
	exist, err := kubeClient.IsDevstreamNS(opts.Namespace)
	if err != nil {
		log.Debugf("Failed to check whether namespace: %s exists.", opts.Namespace)
		return err
	}

	if exist {
		log.Debug("Prepare to delete namespace.")
		if err = kubeClient.DeleteNamespace(opts.Namespace); err != nil {
			if !strings.Contains(err.Error(), "not found") {
				return err
			}
		}
	}

	return nil
}

// Delete application by goclient
func DeleteApp(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	// 1. Create k8s clientset
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	// 2. Delete application
	if err = deleteApp(kubeClient, &opts); err != nil {
		return err
	}

	return nil
}

// GetStatus checks plugin status by goclient
func GetStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	opts, err := NewOptions(options)
	if err != nil {
		return nil, err
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	ready, err := checkDeploymentsAndServicesReady(kubeClient, &opts)
	if err != nil {
		return nil, err
	}

	if !ready {
		return statemanager.ResourceStatus{
			"stopped": true,
		}, nil
	}

	return statemanager.ResourceStatus{
		"running": true,
	}, nil
}
