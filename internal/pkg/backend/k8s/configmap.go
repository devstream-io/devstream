package k8s

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"

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
	// apply configmap
	configMapRes, err := b.client.ApplyConfigMap(b.configMapName, b.namespace, data, labels)
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
	configMap, err := b.client.GetConfigMap(b.configMapName, b.namespace)
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
