package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (b *Backend) applyConfigMap(content string) (*v1.ConfigMap, error) {
	if err := b.client.UpsertNameSpace(b.namespace); err != nil {
		return nil, err
	}

	// build configmap object
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

	// apply configmap
	configMapRes, err := b.client.CoreV1().ConfigMaps(b.namespace).Apply(context.Background(), configMap, applyOptions)
	if err != nil {
		return nil, err
	}
	log.Debugf("configmap %s created, detail: %s", configMapRes.Name, configMapRes.String())

	return configMapRes, nil
}

func (b *Backend) getOrCreateConfigMap() (*v1.ConfigMap, error) {
	configMap, exist, err := b.getConfigMap()
	if err != nil {
		return nil, err
	}
	if exist {
		return configMap, nil
	}

	// if configmap not exist, create it
	log.Infof("configmap %s in namespace %s not exist, will create it", b.configMapName, b.namespace)
	return b.applyConfigMap("")
}

func (b *Backend) getConfigMap() (cm *v1.ConfigMap, exist bool, err error) {
	configMap, err := b.client.CoreV1().ConfigMaps(b.namespace).Get(context.Background(), b.configMapName, metav1.GetOptions{})
	// if configmap not exist, return nil
	if errors.IsNotFound(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	log.Debugf("configmap %s in namespace %s found, detail: %v", configMap.Name, configMap.Namespace, configMap)

	return configMap, true, nil
}
