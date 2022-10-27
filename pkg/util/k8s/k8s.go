package k8s

import (
	"fmt"
	"os"
	"path/filepath"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type Client struct {
	clientset kubernetes.Interface
	// maybe it is not proper to put argocd client in the "k8s client"
	argocd *argocdclient.Clientset
}

var fakeClient *Client

func NewClient() (*Client, error) {
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
