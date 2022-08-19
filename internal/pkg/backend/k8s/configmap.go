package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/tools/record/util"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func (b *Backend) applyConfigMap(content string) (*v1.ConfigMap, error) {
	if err := b.createNamespaceIfNotExist(); err != nil {
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
	log.Info("configmap not exist, will create it")
	return b.applyConfigMap("")
}

func (b *Backend) getConfigMap() (cm *v1.ConfigMap, exist bool, err error) {
	configMap, err := b.client.CoreV1().ConfigMaps(b.namespace).Get(context.Background(), b.configMapName, metav1.GetOptions{})
	// if configmap not exist, return nil
	if util.IsKeyNotFoundError(err) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	log.Debugf("configmap %s found, detail: %s", configMap.Name, configMap.String())

	return configMap, true, nil
}

func (b *Backend) createNamespaceIfNotExist() error {
	_, err := b.client.CoreV1().Namespaces().Get(context.Background(), b.namespace, metav1.GetOptions{})
	// if namespace not exist, try to create it
	if util.IsKeyNotFoundError(err) {
		log.Infof("namespace %s not exist, will create it", b.namespace)
		_, err = b.client.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: b.namespace,
			},
		}, metav1.CreateOptions{})
	}

	return err
}
