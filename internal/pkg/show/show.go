package show

import (
	"github.com/devstream-io/devstream/internal/pkg/show/config"
	"github.com/devstream-io/devstream/internal/pkg/show/status"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type Info string

const (
	ConfigInfo Info = "config"
	StatusInfo Info = "status"
)

var InfoSet = map[Info]struct{}{
	ConfigInfo: {},
	StatusInfo: {},
}

func IsValideInfo(info Info) bool {
	_, ok := InfoSet[info]
	return ok
}

func GenerateInfo(info Info) error {
	switch info {
	case ConfigInfo:
		log.Debugf("Info: %s.", info)
		return config.Show()
	case StatusInfo:
		log.Debugf("Info: %s.", info)
		return status.Show()
	default:
		panic("This should be never happen!")
	}
}
