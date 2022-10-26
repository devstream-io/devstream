package configGetter

import (
	"fmt"
	"os"
)

type EnvGetter struct {
	key string
}

func NewEnvGetter(key string) ItemGetter {
	return &EnvGetter{key: key}
}

func (g *EnvGetter) Get() string {
	return os.Getenv(g.key)
}

func (g *EnvGetter) DescribeWhereToSet() string {
	return fmt.Sprintf("<%s> in environment variable", g.key)
}
