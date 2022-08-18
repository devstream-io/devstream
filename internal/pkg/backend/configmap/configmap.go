package configmap

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/tools/record/util"

	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func (b *ConfigMap) applyConfigMap(content string) (*v1.ConfigMap, error) {
	if err := createNamespaceIfNotExist(b.namespace); err != nil {
		return nil, err
	}

	client, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	labels := map[string]string{
		"app.kubernetes.io/name":       b.configMapName,
		"app.kubernetes.io/managed-by": "DevStream",
		"created_by":                   "DevStream",
	}
	data := map[string]string{
		stateKey: content,
	}

	configMap := corev1.ConfigMap(b.configMapName, b.namespace).
		WithLabels(labels).
		WithData(data).
		WithImmutable(false)

	applyOptions := metav1.ApplyOptions{
		FieldManager: "DevStream",
	}

	configMapRes, err := client.CoreV1().ConfigMaps(b.namespace).Apply(context.Background(), configMap, applyOptions)
	if err != nil {
		return nil, err
	}
	log.Debugf("configmap %s created, detail: %s", configMapRes.Name, configMapRes.String())

	return configMapRes, nil
}

func (b *ConfigMap) getOrCreateConfigMap() (*v1.ConfigMap, error) {
	configMap, exist, err := b.getConfigMap()
	if err != nil {
		return nil, err
	}
	if exist {
		return configMap, nil
	}

	return b.applyConfigMap("")
}

// getConfigMap returns the configMap content,
func (b *ConfigMap) getConfigMap() (cm *v1.ConfigMap, exist bool, err error) {
	client, err := k8s.NewClient()
	if err != nil {
		return nil, false, err
	}

	configMap, err := client.CoreV1().ConfigMaps(b.namespace).Get(context.Background(), b.configMapName, metav1.GetOptions{})
	if util.IsKeyNotFoundError(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	log.Debugf("configmap %s found, detail: %s", configMap.Name, configMap.String())

	return configMap, true, nil
}

func createNamespaceIfNotExist(namespace string) error {
	client, err := k8s.NewClient()
	if err != nil {
		return err
	}

	_, err = client.CoreV1().Namespaces().Get(context.Background(), namespace, metav1.GetOptions{})
	if util.IsKeyNotFoundError(err) {
		_, err = client.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
