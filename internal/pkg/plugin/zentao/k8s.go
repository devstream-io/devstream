package zentao

import (
	"context"
	"errors"
	"strings"

	kerr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DealWithNsWhenInstall(kubeClient *k8s.Client, opts *Options) error {
	// Namespace should not be empty
	if opts.Namespace == "" {
		log.Debugf("No namespace is given.")
		return errors.New("No namespace is given.")
	}
	log.Debugf("Prepare to create the namespace: %s.", opts.Namespace)

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
// If zentaoNs=false, keep namespace and delete other installed resources
func DealWithErrWhenInstall(kubeClient *k8s.Client, opts *Options, step int) error {
	exist, err := IsDevstreamNSExists(kubeClient, opts.Namespace)
	if err != nil {
		log.Debugf("Failed to check whether namespace: %s exists.", opts.Namespace)
		return err
	}

	// zentao namespace is controlled by dtm, just delete this namespace
	if !exist {
		log.Debugf("Prepare to delete the namespace: %s.", opts.Namespace)

		err := kubeClient.DeleteNamespace(opts.Namespace)
		if err != nil {
			log.Debugf("Failed to delete the namespace: %s.", opts.Namespace)
			return err
		}

		log.Debugf("The namespace %s has been deleted.", opts.Namespace)
		return nil
	}

	// Delete installed resources according to `step`
	for i := step; i >= 0; i-- {
		switch installStep[i] {
		case "pv":
			log.Debug("Prepare to delete the persistentvolume in creation exit.")
			if err := DeletePersistentVolume(kubeClient, opts); err != nil {
				if !strings.Contains(err.Error(), "not found") {
					return err
				}
			}
		case "pvc":
			log.Debug("Prepare to delete the persistentvolumeclaim in creation exit.")
			if err := DeletePersistentVolumeClaim(kubeClient, opts); err != nil {
				if !strings.Contains(err.Error(), "not found") {
					return err
				}
			}
		case "app":
			log.Debug("Prepare to delete the application in creation exit.")
			if err := DeleteZentaoAPP(kubeClient, opts); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create the pv zentao needs: zentaoPV and mysqlPV
// Now use `HostPath` and the path has been specified by the official image
func CreatePersistentVolume(kubeClient *k8s.Client, opts *Options) error {

	zentaoPVOption := &k8s.PVOption{
		Name:             opts.PersistentVolume.ZentaoPVName,
		StorageClassName: opts.StorageClassName,
		AccessMode: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Capacity:                      opts.PersistentVolume.ZentaoPVCapacity,
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
		HostPath:                      "/www/zentaopms",
	}

	mysqlPVOption := &k8s.PVOption{
		Name:             opts.PersistentVolume.MysqlPVName,
		StorageClassName: opts.StorageClassName,
		AccessMode: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Capacity:                      opts.PersistentVolume.MysqlPVCapacity,
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
		HostPath:                      "/var/lib/mysql",
	}

	log.Debugf("Prepare to create persistentVolume: %s.", zentaoPVOption.Name)
	if err := kubeClient.CreatePersistentVolume(zentaoPVOption); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The PersistentVolume %s is already exists.", zentaoPVOption.Name)
	}
	log.Debugf("The PersistentVolume %s has been created.", zentaoPVOption.Name)

	log.Debugf("Prepare to create persistentVolume: %s.", mysqlPVOption.Name)
	if err := kubeClient.CreatePersistentVolume(mysqlPVOption); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The PersistentVolume %s is already exists.", mysqlPVOption.Name)
	}
	log.Debugf("The PersistentVolume %s has been created.", mysqlPVOption.Name)

	return nil
}

// Create the pvc zentao needs: zentaoPVC and mysqlPVC
func CreatePersistentVolumeClaim(kubeClient *k8s.Client, opts *Options) error {

	zentaoPVCOption := &k8s.PVCOption{
		Name:             opts.PersistentVolumeClaim.ZentaoPVCName,
		NameSpace:        opts.Namespace,
		StorageClassName: opts.StorageClassName,
		AccessMode: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Requirement: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(opts.PersistentVolumeClaim.ZentaoPVCCapacity),
			},
		},
	}

	mysqlPVCOption := &k8s.PVCOption{
		Name:             opts.PersistentVolumeClaim.MysqlPVCName,
		NameSpace:        opts.Namespace,
		StorageClassName: opts.StorageClassName,
		AccessMode: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Requirement: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceName(corev1.ResourceStorage): resource.MustParse(opts.PersistentVolumeClaim.MysqlPVCCapacity),
			},
		},
	}

	log.Debugf("Prepare to create persistentVolumeClaim: %s.", zentaoPVCOption.Name)
	if err := kubeClient.CreatePersistentVolumeClaim(zentaoPVCOption); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The PersistentVolumeClaim %s is already exists.", zentaoPVCOption.Name)
	}
	log.Debugf("The PersistentVolumeClaim %s has been created.", zentaoPVCOption.Name)

	log.Debugf("Prepare to create persistentVolumeClaim: %s.", mysqlPVCOption.Name)
	if err := kubeClient.CreatePersistentVolumeClaim(mysqlPVCOption); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The PersistentVolumeClaim %s is already exists.", mysqlPVCOption.Name)
	}
	log.Debugf("The PersistentVolumeClaim %s has been created.", mysqlPVCOption.Name)

	return nil
}

func CreateDeployment(kubeClient *k8s.Client, opts *Options) error {

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: opts.Deployment.Name,
			Labels: map[string]string{
				"app": "zentao",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32Ptr(int32(opts.Deployment.Replicas)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "zentao",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "zentao",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "zentao",
							Image: opts.Deployment.Image,
							Env: []corev1.EnvVar{
								{
									Name:  opts.Deployment.MysqlPasswdName,
									Value: opts.Deployment.MysqlPasswdValue,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "zentao",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
								{
									Name:          "mysql",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 3306,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: opts.PersistentVolumeClaim.ZentaoPVCName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: opts.PersistentVolumeClaim.ZentaoPVCName,
								},
							},
						},
						{
							Name: opts.PersistentVolumeClaim.MysqlPVCName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: opts.PersistentVolumeClaim.MysqlPVCName,
								},
							},
						},
					},
				},
			},
		},
	}

	log.Debugf("The Deployment %s has been created.", deployment.Name)
	if _, err := kubeClient.AppsV1().Deployments(opts.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{}); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The Deployment %s is already exists.", deployment.Name)
	}
	log.Debugf("The Deployment %s has been created.", deployment.Name)

	return nil
}

func CreateService(kubeClient *k8s.Client, opts *Options) error {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: opts.Service.Name,
			Labels: map[string]string{
				"app": "zentao",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "zentao",
					Port:       80,
					TargetPort: intstr.IntOrString{IntVal: 80},
					NodePort:   int32(opts.Service.NodePort),
				},
			},
			Selector: map[string]string{
				"app": "zentao",
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	log.Debugf("The Service %s has been created.", svc.Name)
	if _, err := kubeClient.CoreV1().Services(opts.Namespace).Create(context.TODO(), svc, metav1.CreateOptions{}); err != nil {
		if !kerr.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The Service %s is already exists.", svc.Name)
	}
	log.Debugf("The Service %s has been created.", svc.Name)

	return nil
}

// Check whether zentao deployment is ready and zentao service exists
func CheckDeploymentsAndServicesReady(kubeClient *k8s.Client, opts *Options) (bool, error) {
	namespace := opts.Namespace
	deploy := opts.Deployment.Name
	svc := opts.Service.Name

	dp, err := kubeClient.GetDeployment(namespace, deploy)
	if err != nil {
		log.Debugf("Get zentao deployment err: %s", err.Error())
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
		log.Debugf("Get zentao service err: %s", err.Error())
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	log.Debugf("The service %s is ready.", svc)

	return true, nil
}

func DeleteService(kubeClient *k8s.Client, opts *Options) error {
	return kubeClient.CoreV1().Services(opts.Namespace).
		Delete(context.TODO(), opts.Service.Name, metav1.DeleteOptions{})
}

func DeleteDeployment(kubeClient *k8s.Client, opts *Options) error {
	return kubeClient.AppsV1().Deployments(opts.Namespace).
		Delete(context.TODO(), opts.Deployment.Name, metav1.DeleteOptions{})
}

func DeletePersistentVolumeClaim(kubeClient *k8s.Client, opts *Options) error {
	if err := kubeClient.CoreV1().PersistentVolumeClaims(opts.Namespace).
		Delete(context.TODO(), opts.PersistentVolumeClaim.ZentaoPVCName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return kubeClient.CoreV1().PersistentVolumeClaims(opts.Namespace).
		Delete(context.TODO(), opts.PersistentVolumeClaim.MysqlPVCName, metav1.DeleteOptions{})
}

func DeletePersistentVolume(kubeClient *k8s.Client, opts *Options) error {
	if err := kubeClient.CoreV1().PersistentVolumes().
		Delete(context.TODO(), opts.PersistentVolume.ZentaoPVName, metav1.DeleteOptions{}); err != nil {
		return err
	}
	return kubeClient.CoreV1().PersistentVolumes().
		Delete(context.TODO(), opts.PersistentVolume.MysqlPVName, metav1.DeleteOptions{})
}

func DeleteNamespace(kubeClient *k8s.Client, opts *Options) error {
	return kubeClient.CoreV1().Namespaces().Delete(context.TODO(), opts.Namespace, metav1.DeleteOptions{})
}

func getDevstreamNamespace(kubeClient *k8s.Client) (*corev1.NamespaceList, error) {
	return kubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{LabelSelector: "created_by=DevStream"})
}

// Check whether the given namespace is created by dtm
// If the given namespace has label "created_by=DevStream", we'll control it.
// 1. The specified namespace is created by zentao plugin, then it should be deleted
//    when errors are encountered during creation or `dtm delete`.
// 2. The specified namespace is controlled by user, maybe they want to deploy zentao in
//    an existing namespace or other situations, then we should not delete this namespace.
func IsDevstreamNSExists(kubeClient *k8s.Client, namespace string) (bool, error) {
	nsList, err := getDevstreamNamespace(kubeClient)
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
