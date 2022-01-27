package k8s

import (
	"fmt"
	"path/filepath"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	*kubernetes.Clientset
	Argocd *argocdclient.Clientset
}

func NewClient() (*Client, error) {
	homePath := homedir.HomeDir()
	if homePath == "" {
		return nil, fmt.Errorf("failed to get the home directory")
	}

	kubeconfig := filepath.Join(homePath, ".kube", "config")

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
		Clientset: clientset,
		Argocd:    argocdClientset,
	}, nil
}
