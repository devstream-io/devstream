package configGetter

import (
	"fmt"
	"os"
)

type envGetter struct {
	key string
}

func NewEnvGetter(key string) ItemGetter {
	return &envGetter{key: key}
}

func (g *envGetter) Get() string {
	return os.Getenv(g.key)
}

func (g *envGetter) DescribeWhereToSet() string {
	return fmt.Sprintf("<%s> in environment variable", g.key)
}
