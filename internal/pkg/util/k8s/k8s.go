package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	*kubernetes.Clientset
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
	return &Client{
		Clientset: clientset,
	}, nil
}

func (c *Client) ListDeployments(namespace string) ([]appsv1.Deployment, error) {
	dpList, err := c.Clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return dpList.Items, nil
}

func (c *Client) IsDeploymentReady(deployment *appsv1.Deployment) bool {
	return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
}
