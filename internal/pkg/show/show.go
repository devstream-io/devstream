package show

import (
	"github.com/devstream-io/devstream/internal/pkg/show/config"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type Info string

const (
	ConfigInfo Info = "config"
)

var InfoSet = map[Info]struct{}{
	ConfigInfo: {},
}

func IsValideInfo(info Info) bool {
	_, ok := InfoSet[info]
	return ok
}

func GenerateInfo(info Info) error {
	switch info {
	case ConfigInfo:
		log.Debugf("Info: %s.", ConfigInfo)
		return config.Show()
	default:
		panic("This should be never happen!")
	}
}
