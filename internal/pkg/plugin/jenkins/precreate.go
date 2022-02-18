package jenkins

import (
	"os"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/merico-dev/stream/internal/pkg/log"
	"github.com/merico-dev/stream/pkg/util/k8s"
)

const (
	JenkinsName          = "jenkins"
	JenkinsNamespace     = "jenkins"
	JenkinsDataDirectory = "/data/jenkins-volume/"
	JenkinsUid           = 1000
	JenkinsGid           = 1000
)

// See the docs below for more info:
// https://www.jenkins.io/doc/book/installing/kubernetes/
// https://raw.githubusercontent.com/jenkins-infra/jenkins.io/master/content/doc/tutorials/kubernetes/installing-jenkins-on-kubernetes/jenkins-volume.yaml
// https://raw.githubusercontent.com/jenkins-infra/jenkins.io/master/content/doc/tutorials/kubernetes/installing-jenkins-on-kubernetes/jenkins-sa.yaml
func preCreate() error {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}

	if err = initDataDirectory(); err != nil {
		return err
	}

	if err = createPersistentVolume(kubeClient); err != nil {
		return err
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
	pvOption := &k8s.PVOption{
		Name:             "jenkins-pv",
		StorageClassName: "jenkins-pv",
		AccessMode: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Capacity:                      "10Gi",
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
		HostPath:                      "/data/jenkins-volume/",
	}

	if err := kubeClient.CreatePersistentVolume(pvOption); err != nil {
		return err
	}

	return nil
}

func createServiceAccount(kubeClient *k8s.Client) error {
	if err := kubeClient.CreateServiceAccount(JenkinsName, JenkinsNamespace); err != nil {
		return err
	}

	return nil
}

func createClusterRole(kubeClient *k8s.Client) error {
	crOption := &k8s.CROption{
		Name: JenkinsName,
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
		return err
	}

	return nil
}

func createClusterRoleBinding(kubeClient *k8s.Client) error {
	crbOption := &k8s.CRBOption{
		Name:    JenkinsName,
		SANames: []string{JenkinsName},
		RName:   JenkinsName,
	}

	if err := kubeClient.CreateClusterRoleBinding(crbOption); err != nil {
		return err
	}

	return nil
}

func initDataDirectory() error {
	if err := os.MkdirAll(JenkinsDataDirectory, 0755); err != nil {
		log.Errorf("Faield to create data directory: %s", err)
		return err
	}
	log.Debugf("The data directory is created: %s", JenkinsDataDirectory)

	if err := os.Chown(JenkinsDataDirectory, JenkinsUid, JenkinsGid); err != nil {
		log.Errorf("Failed to change the permissions with %s to allow the jenkins account to write its data. %s",
			JenkinsDataDirectory, err)
		return err
	}

	return nil
}
