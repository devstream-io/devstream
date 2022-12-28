package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	argocdv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	argocdclient "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type K8sAPI interface {
	// secret API
	GetSecret(namespace, name string) (map[string]string, error)
	ApplySecret(name, namespace string, data map[string][]byte, labels map[string]string) (*corev1.Secret, error)
	// service API
	CreateService(namespace string, service *corev1.Service) error
	DeleteService(namespace, serviceName string) error
	GetService(namespace, name string) (*corev1.Service, error)
	// storage API
	CreatePersistentVolume(option *PVOption) error
	DeletePersistentVolume(pvName string) error
	CreatePersistentVolumeClaim(opt *PVCOption) error
	DeletePersistentVolumeClaim(namespace, pvcName string) error
	// resource API
	GetResourceStatus(nameSpace string, anFilter, labelFilter map[string]string) (*AllResourceStatus, error)
	ListDeploymentsWithLabel(namespace string, labelFilter map[string]string) ([]appsv1.Deployment, error)
	GetDeployment(namespace, name string) (*appsv1.Deployment, error)
	CreateDeployment(namespace string, deployment *appsv1.Deployment) error
	WaitForDeploymentReady(retry int, namespace, deployName string) error
	DeleteDeployment(namespace, deployName string) error
	ListDaemonsetsWithLabel(namespace string, labeFilter map[string]string) ([]appsv1.DaemonSet, error)
	GetStatefulset(namespace, name string) (*appsv1.StatefulSet, error)
	// namespace API
	UpsertNameSpace(nameSpace string) error
	GetNamespace(namespace string) (*corev1.Namespace, error)
	IsDevstreamNS(namespace string) (bool, error)
	CreateNamespace(namespace string) error
	DeleteNamespace(namespace string) error
	IsNamespaceExists(namespace string) (bool, error)
	// configmap API
	ApplyConfigMap(name, namespace string, data, labels map[string]string) (*corev1.ConfigMap, error)
	GetConfigMap(name, namespace string) (*corev1.ConfigMap, error)
	// argocd API
	ListArgocdApplications(namespace string) ([]argocdv1alpha1.Application, error)
	GetArgocdApplication(namespace, name string) (*argocdv1alpha1.Application, error)
	IsArgocdApplicationReady(application *argocdv1alpha1.Application) bool
	DescribeArgocdApp(app *argocdv1alpha1.Application) map[string]interface{}
}

type Client struct {
	clientset kubernetes.Interface
	// maybe it is not proper to put argocd client in the "k8s client"
	argocd *argocdclient.Clientset
}

var fakeClient *Client

func NewClient() (K8sAPI, error) {
	// if UseFakeClient() is called, return the fake client.
	if fakeClient != nil {
		return fakeClient, nil
	}

	// TL;DR: Don't use viper.GetString("xxx") in the `util/xxx` package.
	// Don't use `kubeconfig := viper.GetString("kubeconfig")` here,
	// it will fail without calling `viper.BindEnv("github_token")` first.
	// os.Getenv() function is more clear and reasonable here.
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("kubeconfig")
	}
	if kubeconfig != "" {
		log.Debugf("Got the kubeconfig from env: %s.", kubeconfig)
	} else {
		log.Debugf("Failed to get the kubecondig from env. Prepare to get it from home dir.")
		homePath := homedir.HomeDir()
		if homePath == "" {
			return nil, fmt.Errorf("failed to get the home directory")
		}

		kubeconfig = filepath.Join(homePath, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	argocdClientset, err := argocdclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientset,
		argocd:    argocdClientset,
	}, nil
}

// UseFakeClient is used for testing,
// if this function is called, NewClient() will return the fake client.
func UseFakeClient(k8sClient kubernetes.Interface, argoClient *argocdclient.Clientset) {
	fakeClient = &Client{
		clientset: k8sClient,
		argocd:    argoClient,
	}
}
