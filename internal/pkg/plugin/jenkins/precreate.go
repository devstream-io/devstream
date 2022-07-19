package jenkins

import (
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/homedir"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	jenkinsName                      = "jenkins"
	jenkinsNamespace                 = "jenkins"
	jenkinsDataDirectory             = "/data/jenkins-volume/"
	jenkinsPvName                    = "jenkins-pv"
	jenkinsPvDefaultStorageClassName = "jenkins-pv"
)

// See the docs below for more info:
// https://www.jenkins.io/doc/book/installing/kubernetes/
// https://raw.githubusercontent.com/jenkins-infra/jenkins.io/master/content/doc/tutorials/kubernetes/installing-jenkins-on-kubernetes/jenkins-volume.yaml
// https://raw.githubusercontent.com/jenkins-infra/jenkins.io/master/content/doc/tutorials/kubernetes/installing-jenkins-on-kubernetes/jenkins-sa.yaml
func preCreate(options plugininstaller.RawOptions) error {
	opts, err := newOptions(options)
	if err != nil {
		return err
	}
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	if opts.TestEnv {
		log.Info("Test environment is enabled. Please ensure you have created the directories correctly under the guide of plugin doc.")
		if err = createPersistentVolume(kubeClient); err != nil {
			return err
		}
	}

	if err = createServiceAccount(kubeClient); err != nil {
		return err
	}

	if err = createClusterRole(kubeClient); err != nil {
		return err
	}

	if err = createClusterRoleBinding(kubeClient); err != nil {
		return err
	}

	return nil
}

func createPersistentVolume(kubeClient *k8s.Client) error {
	dataDir := getRealJenkinsDataDirectory()
	if dataDir == "" {
		return fmt.Errorf("failed to get the real Jenkins data directory")
	}

	pvOption := &k8s.PVOption{
		Name:             jenkinsPvName,
		StorageClassName: jenkinsPvDefaultStorageClassName,
		AccessMode: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Capacity:                      "20Gi",
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
		HostPath:                      dataDir,
	}

	if err := kubeClient.CreatePersistentVolume(pvOption); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The resource %s is already exists.", "PersistentVolume")
	}

	return nil
}

func createServiceAccount(kubeClient *k8s.Client) error {
	if err := kubeClient.CreateServiceAccount(jenkinsName, jenkinsNamespace); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The resource %s is already exists.", "ServiceAccount")
	}

	return nil
}

func createClusterRole(kubeClient *k8s.Client) error {
	crOption := &k8s.CROption{
		Name: jenkinsName,
		PolicyRules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{
					"statefulsets",
					"services",
					"replicationcontrollers",
					"replicasets",
					"podtemplates",
					"podsecuritypolicies",
					"pods",
					"pods/log",
					"pods/exec",
					"podpreset",
					"poddisruptionbudget",
					"persistentvolumes",
					"persistentvolumeclaims",
					"jobs",
					"endpoints",
					"deployments",
					"deployments/scale",
					"daemonsets",
					"cronjobs",
					"configmaps",
					"namespaces",
					"events",
					"secrets",
				},
				Verbs: []string{
					"create",
					"get",
					"watch",
					"delete",
					"list",
					"patch",
					"update",
				},
			},
			{
				APIGroups: []string{""},
				Resources: []string{
					"nodes",
				},
				Verbs: []string{
					"get",
					"watch",
					"list",
					"update",
				},
			},
		},
	}

	if err := kubeClient.CreateClusterRole(crOption); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The resource %s is already exists.", "ClusterRole")
	}

	return nil
}

func createClusterRoleBinding(kubeClient *k8s.Client) error {
	crbOption := &k8s.CRBOption{
		Name:    jenkinsName,
		SANames: []string{jenkinsName},
		RName:   jenkinsName,
	}

	if err := kubeClient.CreateClusterRoleBinding(crbOption); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The resource %s is already exists.", "ClusterRoleBinding")
	}

	return nil
}

// get the data directory of Jenkins by the home dir of the machine where dtm is running
func getRealJenkinsDataDirectory() string {
	home := homedir.HomeDir()
	if home == "" {
		log.Errorf("Failed to get the homedir.")
		return ""
	}

	log.Debugf("Got the homedir: %s", home)
	realPath := filepath.Join(home, jenkinsDataDirectory)
	log.Debugf("The real Jenkins data directory is: %s", realPath)
	return realPath
}
