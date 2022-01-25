package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (c *Client) ListDeployments(namespace string) ([]appsv1.Deployment, error) {
	dpList, err := c.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return dpList.Items, nil
}

func (c *Client) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	return c.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *Client) IsDeploymentReady(deployment *appsv1.Deployment) bool {
	return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas
}

func (c *Client) ListDaemonsets(namespace string) ([]appsv1.DaemonSet, error) {
	dsList, err := c.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return dsList.Items, nil
}

func (c *Client) GetDaemonset(namespace, name string) (*appsv1.DaemonSet, error) {
	return c.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *Client) IsDaemonsetReady(daemonset *appsv1.DaemonSet) bool {
	return daemonset.Status.NumberReady == daemonset.Status.DesiredNumberScheduled
}

func (c *Client) ListStatefulsets(namespace string) ([]appsv1.StatefulSet, error) {
	ssList, err := c.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ssList.Items, nil
}

func (c *Client) GetStatefulset(namespace, name string) (*appsv1.StatefulSet, error) {
	return c.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *Client) IsStatefulsetReady(statefulset *appsv1.StatefulSet) bool {
	return statefulset.Status.ReadyReplicas == *statefulset.Spec.Replicas
}

func (c *Client) ListServices(namespace string) ([]corev1.Service, error) {
	services, err := c.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return services.Items, nil
}

func (c *Client) GetService(namespace, name string) (*corev1.Service, error) {
	return c.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
