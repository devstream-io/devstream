package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

func (c *Client) GetNamespace(namespace string) (*corev1.Namespace, error) {
	return c.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
}

func (c *Client) CreateNamespace(namespace string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"created_by": "DevStream",
			},
		},
	}
	_, err := c.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	return err
}

func (c *Client) DeleteNamespace(namespace string) error {
	if namespace == "default" || namespace == "kube-system" {
		return fmt.Errorf("you can't delete the default or kube-system namespace")
	}

	gracePeriodSeconds := int64(0)
	return c.CoreV1().Namespaces().Delete(
		context.TODO(), namespace, metav1.DeleteOptions{GracePeriodSeconds: &gracePeriodSeconds})
}

func (c *Client) IsNamespaceExists(namespace string) (bool, error) {
	_, err := c.GetNamespace(namespace)
	if err != nil && !errors.IsNotFound(err) {
		return false, err
	}
	// not exist
	if err != nil {
		return false, nil
	}
	// exist
	return true, nil
}
