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

func (c *Client) ListDevstreamNamespace() (*corev1.NamespaceList, error) {
	return c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{LabelSelector: "created_by=DevStream"})
}

// Check whether the given namespace is created by dtm
// If the given namespace has label "created_by=DevStream", we'll control it.
// 1. The specified namespace is created by dtm, then it should be deleted
//    when errors are encountered during creation or `dtm delete`.
// 2. The specified namespace is controlled by user, maybe they want to deploy plugins in
//    an existing namespace or other situations, then we should not delete this namespace.
func (c *Client) IsDevstreamNS(namespace string) (bool, error) {
	nsList, err := c.ListDevstreamNamespace()
	if err != nil {
		// not exist
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	for _, ns := range nsList.Items {
		// exist
		if ns.ObjectMeta.Name == namespace {
			return true, nil
		}
	}
	return false, nil
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

func (c *Client) DeleteService(namespace, serviceName string) error {
	return c.CoreV1().Services(namespace).
		Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
}
