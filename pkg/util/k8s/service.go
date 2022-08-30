package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (c *Client) GetService(namespace, name string) (*corev1.Service, error) {
	return c.clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (c *Client) DeleteService(namespace, serviceName string) error {
	return c.clientset.CoreV1().Services(namespace).
		Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
}

func (c *Client) CreateService(namespace string, service *corev1.Service) error {
	_, err := c.clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if kerr.IsAlreadyExists(err) {
		log.Infof("The Service %s is already exists.", service.Name)
		return nil
	}
	return err
}
