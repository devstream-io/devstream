package backend

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type Type string

const Local Type = "local"

// Backend is used to persist data, it can be local file/etcd/s3/...
type Backend interface {
	// Read is used to reads data from persistent storage.
	Read() ([]byte, error)
	// Write is used to writes data to persistent storage.
	Write(data []byte) error
}

// GetBackend will return a Backend by the given name.
func GetBackend(typeName Type) (Backend, error) {
	switch typeName {
	case Local:
		log.Debugf("Used the Backend: %s.", typeName)
		return local.NewLocal(local.DefaultStateFile), nil
	default:
		return nil, fmt.Errorf("the backend type < %s > is illegal", typeName)
	}
}
