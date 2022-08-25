package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
)

func (c *Client) ApplyConfigMap(name, namespace string, data, labels map[string]string) (*v1.ConfigMap, error) {
	configMap := corev1.ConfigMap(name, namespace).
		WithLabels(labels).
		WithData(data).
		WithImmutable(false)
	applyOptions := metav1.ApplyOptions{
		FieldManager: "DevStream",
	}
	return c.clientset.CoreV1().ConfigMaps(namespace).Apply(context.Background(), configMap, applyOptions)
}

func (c *Client) GetConfigMap(name, namespace string) (*v1.ConfigMap, error) {
	return c.clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
}
