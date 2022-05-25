package configloader

// OriginalConfig is the config YAML user writes
type OriginalConfig struct {
	VarFile  string `yaml:"varFile"`
	ToolFile string `yaml:"toolFile"`
	State    *State `yaml:"state"`
}

type State struct {
	Backend string             `yaml:"backend"`
	Options StateConfigOptions `yaml:"options"`
}

type StateConfigOptions struct {
	// for s3 backend
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
	Key    string `yaml:"key"`
	// for local backend
	StateFile string `yaml:"stateFile"`
}
