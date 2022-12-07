package configmanager

import (
	"github.com/imdario/mergo"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
)

// appSpec is app special options
type appSpec struct {
	// language config
	Language  string `yaml:"language" mapstructure:"language"`
	FrameWork string `yaml:"framework" mapstructure:"framework"`
}

// merge will merge vars and appSpec
func (s *appSpec) merge(vars map[string]any) map[string]any {
	specMap, err := mapz.DecodeStructToMap(s)
	if err != nil {
		log.Warnf("appspec %+v decode failed: %+v", s, err)
		return map[string]any{}
	}
	_ = mergo.Merge(&specMap, vars)
	return specMap
}

func (s *appSpec) updatePiplineOption(options RawOptions) {
	if _, exist := options["language"]; !exist && s.hasLanguageConfig() {
		options["language"] = RawOptions{
			"name":      s.Language,
			"framework": s.FrameWork,
		}
	}
}

func (s *appSpec) hasLanguageConfig() bool {
	return s.Language != "" || s.FrameWork != ""
}
