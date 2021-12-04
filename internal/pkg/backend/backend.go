package backend

import (
	"fmt"

	"github.com/merico-dev/stream/internal/pkg/backend/local"
)

var backends map[string]Backend

func init() {
	if backends == nil{
		backends = map[string]Backend{
			"local": local.NewLocal(""),
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
func GetBackend(name string) (Backend, error) {
	if b, ok := backends[name]; ok {
		return b, nil
	} else {
		return nil, fmt.Errorf("backend %s is illegal", name)
	}
}
