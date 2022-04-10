package golang

type Options struct {
	Owner     string
	Org       string
	Repo      string
	Branch    string
	ImageRepo string `mapstructure:"image_repo"`
}
