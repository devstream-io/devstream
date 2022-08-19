package backend

import (
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/backend/k8s"
	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/backend/s3"
	"github.com/devstream-io/devstream/internal/pkg/backend/types"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

// Backend is used to persist data, it can be local file/s3/k8s...
type Backend interface {
	// Read is used to read data from persistent storage.
	Read() ([]byte, error)
	// Write is used to write data to persistent storage.
	Write(data []byte) error
}

// GetBackend will return a Backend by the given name.
func GetBackend(state configmanager.State) (Backend, error) {
	typeName := types.Type(strings.ToLower(state.Backend))

	switch typeName {
	case types.Local:
		return local.NewLocal(state.Options.StateFile)
	case types.S3:
		return s3.NewS3Backend(state.Options.Bucket, state.Options.Region, state.Options.Key)
	case types.K8s, types.K8sAlias:
		return k8s.NewBackend(state.Options.Namespace, state.Options.ConfigMap)
	default:
		return nil, types.NewInvalidBackendErr(state.Backend)
	}
}
