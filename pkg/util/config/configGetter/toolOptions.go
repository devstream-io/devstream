package configGetter

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type toolOptionsGetter struct {
	key        string
	rawOptions configmanager.RawOptions
}

func NewToolOptionsGetter(keyPath string, rawOptions configmanager.RawOptions) ItemGetter {
	return &toolOptionsGetter{
		key:        keyPath,
		rawOptions: rawOptions,
	}
}

func (g *toolOptionsGetter) Get() string {
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(g.rawOptions); err != nil {
		log.Warnf("failed to marshal tool options: %v, err: %v", g.rawOptions, err)
		return ""
	}
	tmpViper := viper.New()
	tmpViper.SetConfigType("json")
	if err := tmpViper.ReadConfig(buffer); err != nil {
		log.Warnf("failed to read tool options: %v, err: %v", g.rawOptions, err)
		return ""
	}
	return tmpViper.GetString(g.key)
}

func (g *toolOptionsGetter) DescribeWhereToSet() string {
	return fmt.Sprintf("<%s> in tools.options", g.key)
}
