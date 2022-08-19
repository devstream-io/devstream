package backend

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/backend/k8s"
	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/backend/s3"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type Type string

const (
	Local Type = "local"
	S3    Type = "s3"
	K8s   Type = "k8s"
)

// Backend is used to persist data, it can be local file/etcd/s3/k8s...
type Backend interface {
	// Read is used to read data from persistent storage.
	Read() ([]byte, error)
	// Write is used to write data to persistent storage.
	Write(data []byte) error
}

// GetBackend will return a Backend by the given name.
func GetBackend(state configmanager.State) (Backend, error) {
	typeName := Type(state.Backend)
	switch typeName {
	case Local:
		log.Infof("Using local backend. State file: %s.", state.Options.StateFile)
		return local.NewLocal(state.Options.StateFile), nil
	case S3:
		log.Infof("Using s3 backend. Bucket: %s, region: %s, key: %s.", state.Options.Bucket, state.Options.Region, state.Options.Key)
		return s3.NewS3Backend(
			state.Options.Bucket,
			state.Options.Region,
			state.Options.Key,
		), nil
	case K8s:
		log.Infof("Using configmap backend. Namespace: %s, ConfigMap name: %s.", state.Options.Namespace, state.Options.ConfigMap)
		return k8s.NewBackend(
			state.Options.Namespace,
			state.Options.ConfigMap)
	default:
		return nil, fmt.Errorf("the backend type < %s > is illegal", typeName)
	}
}
