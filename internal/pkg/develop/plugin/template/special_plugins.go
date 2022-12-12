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
	"gitlabci-golang":  NewSpecialPlugin("gitlabci/golang", "golang"),
	"gitlabci-java":    NewSpecialPlugin("gitlabci/java", "java"),
	"gitlabci-generic": NewSpecialPlugin("gitlabci/generic", "generic"),
}
