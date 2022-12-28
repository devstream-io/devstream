package cifile

import (
	"fmt"
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
		renderContent, err := template.NewRenderClient(&template.TemplateOption{
			Name: fmt.Sprintf("ci-%s", ciTemplateName),
		}, template.LocalFileGetter).Render(filePath, vars)
		if err != nil {
			return nil, err
		}
		return []byte(renderContent), nil
	}
}
