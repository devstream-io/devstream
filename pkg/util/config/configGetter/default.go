package configGetter

type defaultGetter struct {
	value string
}

// defaultValue sets the default value for the config item,
// Do not set "" as default value, it will be ignored
func DefaultValue(value string) ItemGetter {
	return &defaultGetter{
		value: value,
	}
}

func (g *defaultGetter) Get() string {
	return g.value
}

func (g *defaultGetter) DescribeWhereToSet() string {
	return ""
}
