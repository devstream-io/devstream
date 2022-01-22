package golang

// TODO(daniel-hutao): Param should keep as same as other plugins named Param or keep as same as plugin github-actions?
type Param struct {
	Owner     string
	Repo      string
	Branch    string
	ImageRepo string `mapstructure:"image_repo"`
}
