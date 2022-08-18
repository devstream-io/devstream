package configmap

import (
	"sync"

	"github.com/devstream-io/devstream/pkg/util/k8s"
)

const (
	defaultNamespace     = "devstream"
	defaultConfigMapName = "devstream-backend"
	stateKey             = "state"
)

type ConfigMap struct {
	mu            sync.Mutex
	namespace     string
	configMapName string
}

// NewBackend returns a ConfigMap instance as backend
func NewBackend(namespace, configMapName string) (*ConfigMap, error) {
	b := &ConfigMap{
		namespace:     namespace,
		configMapName: configMapName,
	}
	if b.namespace == "" {
		b.namespace = defaultNamespace
	}
	if b.configMapName == "" {
		b.configMapName = defaultConfigMapName
	}
	return b, testK8sConnect()
}

func (b *ConfigMap) Read() ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	configMap, err := b.getOrCreateConfigMap()
	if err != nil {
		return nil, err
	}

	return []byte(configMap.Data[stateKey]), nil
}

func (b *ConfigMap) Write(data []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, err := b.applyConfigMap(string(data))
	return err
}

func testK8sConnect() error {
	_, err := k8s.NewClient()
	return err
}
