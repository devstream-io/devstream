package harbordocker

// Options is the struct for configurations of the harbor-docker plugin.
type Options struct {
	// TODO(daniel-hutao): add more options here asap
	Hostname string `validate:"hostname" mapstructure:"hostname"`
}
