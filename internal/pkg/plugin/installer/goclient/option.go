package goclient

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

// Options is the struct for parameters used by the goclient install config.
type Options struct {
	Namespace              string                   `validate:"required"`
	StorageClassName       string                   `validate:"required"`
	PersistentVolumes      []*PersistentVolume      `validate:"required"`
	PersistentVolumeClaims []*PersistentVolumeClaim `validate:"required"`
	Deployment             *Deployment              `validate:"required"`
	Service                *Service                 `validate:"required"`
}

type PersistentVolume struct {
	PVName     string `validate:"required"`
	PVCapacity string `validate:"required"`
	HostPath   string
}

type PersistentVolumeClaim struct {
	PVCName     string `validate:"required"`
	PVCCapacity string `validate:"required"`
}

type Deployment struct {
	Name     string `validate:"required"`
	Replicas int    `validate:"required"`
	Image    string `validate:"required"`
	Envs     []*Env `validate:"required"`
}

type Service struct {
	Name     string `validate:"required"`
	NodePort int    `validate:"required"`
}

type Env struct {
	Key   string `validate:"required"`
	Value string `validate:"required"`
}

// NewOptions create options by raw options
func NewOptions(options configmanager.RawOptions) (Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return opts, err
	}
	return opts, nil
}
