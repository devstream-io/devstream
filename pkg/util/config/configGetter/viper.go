package configGetter

import (
	"fmt"

	"github.com/spf13/viper"
)

// ViperGetter gets value from viper which supports environment variable and command line flag
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
	return fmt.Sprintf("<%s> in environment variable or command line flag", g.key)
}
