package backend

import (
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/backend/k8s"
	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/backend/s3"
	"github.com/devstream-io/devstream/internal/pkg/backend/types"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
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
	typeName := types.Type(state.Backend)
	switch {
	case types.Local == typeName:
		return local.NewLocal(state.Options.StateFile)
	case types.S3 == typeName:
		return s3.NewS3Backend(
			state.Options.Bucket,
			state.Options.Region,
			state.Options.Key)
	case strings.ToLower(state.Backend) == types.K8s.String() ||
		strings.ToLower(state.Backend) == types.K8sAlis.String():
		return k8s.NewBackend(
			state.Options.Namespace,
			state.Options.ConfigMap)
	default:
		return nil, types.NewInvalidBackendErr(state.Backend)
	}
}
