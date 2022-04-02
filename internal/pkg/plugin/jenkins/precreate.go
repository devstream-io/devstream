package jenkins

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/homedir"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
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
	dataDir := getRealJenkinsDataDirectory()
	if dataDir == "" {
		return fmt.Errorf("failed to get the real Jenkins data directory")
	}

	pvOption := &k8s.PVOption{
		Name:             "jenkins-pv",
		StorageClassName: "jenkins-pv",
		AccessMode: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Capacity:                      "10Gi",
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
	if err := kubeClient.CreateServiceAccount(JenkinsName, JenkinsNamespace); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The resource %s is already exists.", "ServiceAccount")
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
		if !errors.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The resource %s is already exists.", "ClusterRole")
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
		if !errors.IsAlreadyExists(err) {
			return err
		}
		log.Infof("The resource %s is already exists.", "ClusterRoleBinding")
	}

	return nil
}

func initDataDirectory() error {
	dataDir := getRealJenkinsDataDirectory()
	if dataDir == "" {
		return fmt.Errorf("failed to get the real Jenkins data directory")
	}

	f, err := os.Stat(dataDir)
	if err != nil && os.IsNotExist(err) {
		return fmt.Errorf("you should create the data directory \"%s\" manually", dataDir)
	}

	if err != nil {
		return fmt.Errorf("failed to stat the data directory \"%s\": %s", dataDir, err)
	}

	uid := int(f.Sys().(*syscall.Stat_t).Uid)
	gid := int(f.Sys().(*syscall.Stat_t).Gid)

	if uid != JenkinsUid || gid != JenkinsGid {
		return fmt.Errorf("you should chown the data directory to 1000:1000. Expected the %s with uid=%d gid=%d, but got the uid=%d gid=%d", dataDir, JenkinsUid, JenkinsGid, uid, gid)
	}

	log.Debugf("The data directory %s is ready.", dataDir)
	return nil
}

func getRealJenkinsDataDirectory() string {
	home := homedir.HomeDir()
	if home == "" {
		log.Errorf("Failed to get the homedir.")
		return ""
	}

	log.Debugf("Got the homedir: %s", home)
	realPath := filepath.Join(home, JenkinsDataDirectory)
	log.Debugf("The real Jenkins data directory is: %s", realPath)
	return realPath
}
