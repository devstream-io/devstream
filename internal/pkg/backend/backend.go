package backend

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/backend/s3"
	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type Type string

const (
	Local Type = "local"
	S3    Type = "s3"
)

// Backend is used to persist data, it can be local file/etcd/s3/...
type Backend interface {
	// Read is used to read data from persistent storage.
	Read() ([]byte, error)
	// Write is used to write data to persistent storage.
	Write(data []byte) error
}

// GetBackend will return a Backend by the given name.
func GetBackend(stateConfig configloader.State) (Backend, error) {
	typeName := Type(stateConfig.Backend)
	switch typeName {
	case Local:
		log.Debugf("Used the Backend: %s.", typeName)
		return local.NewLocal(stateConfig.Options.StateFile), nil
	case S3:
		log.Debugf("Used the Backend: %s.", typeName)
		return s3.NewS3Backend(
			stateConfig.Options.Bucket,
			stateConfig.Options.Region,
			stateConfig.Options.Key,
		), nil
	default:
		return nil, fmt.Errorf("the backend type < %s > is illegal", typeName)
	}
}
