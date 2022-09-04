package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
)

func (c *Client) GetSecret(namespace, name string) (map[string]string, error) {
	secretMap := make(map[string]string)
	secret, err := c.clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	for k, v := range secret.Data {
		// decode with base64
		secretMap[k] = string(v)
	}
	return secretMap, nil
}

func (c *Client) ApplySecret(name, namespace string, data map[string][]byte, labels map[string]string) (*corev1.Secret, error) {
	secret := v1.Secret(name, namespace).
		WithLabels(labels).
		WithData(data).
		WithImmutable(false)

	applyOptions := metav1.ApplyOptions{
		FieldManager: "DevStream",
	}

	return c.clientset.CoreV1().
		Secrets(namespace).
		Apply(context.Background(), secret, applyOptions)
}
