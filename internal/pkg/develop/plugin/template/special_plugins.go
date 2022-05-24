package template

type SpecialPlugin struct {
	DirName     string
	PackageName string
}

func NewSpecialPlugin(dirName string, packageName string) *SpecialPlugin {
	return &SpecialPlugin{
		DirName:     dirName,
		PackageName: packageName,
	}
}

var SpecialPluginNameMap = map[string]*SpecialPlugin{
	"gitlabci-golang":                NewSpecialPlugin("gitlabci/golang", "gitlabci"),
	"gitlabci-generic":               NewSpecialPlugin("gitlabci/generic", "generic"),
	"gitlab-repo-scaffolding-golang": NewSpecialPlugin("reposcaffolding/github/golang", "golang"),
	"github-repo-scaffolding-golang": NewSpecialPlugin("reposcaffolding/gitlab/golang", "golang"),
	"githubactions-golang":           NewSpecialPlugin("githubactions/golang", "golang"),
	"githubactions-nodejs":           NewSpecialPlugin("githubactions/nodejs", "nodejs"),
	"githubactions-python":           NewSpecialPlugin("githubactions/python", "python"),
}
