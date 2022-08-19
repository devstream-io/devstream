package k8s

import (
	"sync"

	"github.com/devstream-io/devstream/pkg/util/k8s"
)

const (
	defaultNamespace     = "devstream"
	defaultConfigMapName = "devstream-state"
	stateKey             = "state"
)

type Backend struct {
	mu sync.Mutex

	namespace     string
	configMapName string

	client *k8s.Client
}

// NewBackend returns a backend which uses ConfigMap to store data
func NewBackend(namespace, configMapName string) (*Backend, error) {
	c, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	b := &Backend{
		namespace:     namespace,
		configMapName: configMapName,
		client:        c,
	}
	if b.namespace == "" {
		b.namespace = defaultNamespace
	}
	if b.configMapName == "" {
		b.configMapName = defaultConfigMapName
	}

	return b, nil
}

func (b *Backend) Read() ([]byte, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	configMap, err := b.getOrCreateConfigMap()
	if err != nil {
		return nil, err
	}

	return []byte(configMap.Data[stateKey]), nil
}

func (b *Backend) Write(data []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	_, err := b.applyConfigMap(string(data))
	return err
}
