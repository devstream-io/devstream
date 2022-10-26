package configGetter

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ViperGetter gets value from viper which supports environment variable(only upper case) and command line flag
type ViperGetter struct {
	key string
}

func NewViperGetter(key string) ItemGetter {
	return &ViperGetter{key: key}
}

func (g *ViperGetter) Get() string {
	return viper.GetString(g.key)
}

func (g *ViperGetter) DescribeWhereToSet() string {
	return fmt.Sprintf("<%s> in environment variable, or <%s> in command line flag", strings.ToUpper(g.key), g.key)
}
