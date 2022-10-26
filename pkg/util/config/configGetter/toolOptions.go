package configGetter

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/spf13/viper"
)

type ToolOptionsGetter struct {
	key        string
	rawOptions configmanager.RawOptions
}

func NewToolOptionsGetter(keyPath string, rawOptions configmanager.RawOptions) ItemGetter {
	return &ToolOptionsGetter{
		key:        keyPath,
		rawOptions: rawOptions,
	}
}

func (g *ToolOptionsGetter) Get() string {
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(g.rawOptions); err != nil {
		log.Warnf("failed to marshal tool options: %v, err: %v", g.rawOptions, err)
		return ""
	}
	tmpViper := viper.New()
	tmpViper.SetConfigType("json")
	tmpViper.ReadConfig(buffer)
	return tmpViper.GetString(g.key)
}

func (g *ToolOptionsGetter) DescribeWhereToSet() string {
	return fmt.Sprintf("<%s> in tools.options", g.key)
}
