package configmanager

type CoreConfig struct {
	State *State `yaml:"state"`
}

// State is the struct for reading the state configuration in the config file.
// It defines how the state is stored, specifies the type of backend and related options.
type State struct {
	Backend string             `yaml:"backend"`
	Options StateConfigOptions `yaml:"options"`
	BaseDir string             `yaml:"-"` // baseDir is the base directory of the config file
}

// StateConfigOptions is the struct for reading the options of the state backend.
type StateConfigOptions struct {
	// for s3 backend
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
	Key    string `yaml:"key"`
	// for local backend
	StateFile string `yaml:"stateFile"`
	// for k8s backend
	Namespace string `yaml:"namespace"`
	ConfigMap string `yaml:"configmap"`
}
