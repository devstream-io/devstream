package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
