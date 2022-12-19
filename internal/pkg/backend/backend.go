package backend

import (
	"fmt"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/backend/k8s"
	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/internal/pkg/backend/s3"
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

// Backend is used to persist data, it can be local file/s3/k8s...
type Backend interface {
	// Read is used to read data from persistent storage.
	Read() ([]byte, error)
	// Write is used to write data to persistent storage.
	Write(data []byte) error
}

type Type string

const (
	LocalBackend    Type = "local"
	S3Backend       Type = "s3"
	K8sBackend      Type = "k8s"
	K8sBackendAlias Type = "kubernetes"
)

// GetBackend will return a Backend by the given name.
func GetBackend(state configmanager.State) (Backend, error) {
	typeName := Type(strings.ToLower(string(state.Backend)))

	switch typeName {
	case LocalBackend:
		return local.NewLocal(state.BaseDir, state.Options.StateFile)
	case S3Backend:
		return s3.NewS3Backend(state.Options.Bucket, state.Options.Region, state.Options.Key)
	case K8sBackend, K8sBackendAlias:
		return k8s.NewBackend(state.Options.Namespace, state.Options.ConfigMap)
	default:
		return nil, fmt.Errorf("the backend type < %s > is illegal", state.Backend)
	}
}
