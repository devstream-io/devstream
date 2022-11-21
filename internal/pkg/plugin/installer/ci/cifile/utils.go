package cifile

import (
	"os"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/template"
)

const (
	ciTemplateName = "ci_template"
)

func processCIFilesFunc(vars map[string]interface{}) file.DirFileContentFunc {
	return func(filePath string) ([]byte, error) {
		if len(vars) == 0 {
			return os.ReadFile(filePath)
		}
		renderContent, err := template.New().FromLocalFile(filePath).SetDefaultRender(ciTemplateName, vars).Render()
		if err != nil {
			return nil, err
		}
		return []byte(renderContent), nil
	}
}
