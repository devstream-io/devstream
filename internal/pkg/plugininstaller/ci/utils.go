package ci

import (
	"os"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/template"
)

func processCIFilesFunc(templateName string, vars map[string]interface{}) file.DirFileProcessFunc {
	return func(filePath string) ([]byte, error) {
		if len(vars) == 0 {
			return os.ReadFile(filePath)
		}
		renderContent, err := template.New().FromLocalFile(filePath).SetDefaultRender(templateName, vars).Render()
		if err != nil {
			return nil, err
		}
		return []byte(renderContent), nil
	}
}
