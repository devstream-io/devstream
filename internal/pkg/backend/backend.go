package backend

import (
	"fmt"

	"github.com/merico-dev/stream/internal/pkg/backend/local"
)

type BackendType string

const BackendLocal BackendType = "local"

var backends map[BackendType]Backend

func init() {
	if backends == nil {
		backends = map[BackendType]Backend{
			BackendLocal: local.NewLocal(local.DefaultStateFile),
		}
	}
}

// Backend is used to persist data, it can be local file/etcd/s3/...
type Backend interface {
	// Read is used to reads data from persistent storage.
	Read() ([]byte, error)
	// Write is used to writes data to persistent storage.
	Write(data []byte) error
}

// GetBackend will return a Backend by the given name.
func GetBackend(name BackendType) (Backend, error) {
	if b, ok := backends[name]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("the backend < %s > is illegal", name)
}
