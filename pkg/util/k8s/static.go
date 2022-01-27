package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func (c *Client) DeleteNamespace(namespace string) error {
	if namespace == "default" || namespace == "kube-system" {
		return fmt.Errorf("you can't delete the default or kube-system namespace")
	}
	return c.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
}
