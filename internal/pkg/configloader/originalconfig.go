package configloader

// OriginalConfig is the struct representing the complete original configuration YAML files.
type OriginalConfig struct {
	VarFile  string `yaml:"varFile"`
	ToolFile string `yaml:"toolFile"`
	State    *State `yaml:"state"`
}

// State is the struct for reading the state configuration in the config file.
// It defines how the state is stored, specifies the type of backend and related options.
type State struct {
	Backend string             `yaml:"backend"`
	Options StateConfigOptions `yaml:"options"`
}

// StateConfigOptions is the struct for reading the options of the state backend.
type StateConfigOptions struct {
	// for s3 backend
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
	Key    string `yaml:"key"`
	// for local backend
	StateFile string `yaml:"stateFile"`
}
