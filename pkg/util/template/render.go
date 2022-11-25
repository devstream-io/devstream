package template

import (
	"bytes"
	"os"
	"text/template"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Render(name, templateStr string, variable any, funcMaps ...template.FuncMap) (string, error) {
	t := template.New(name).Option("missingkey=error").Delims("[[", "]]")

	for _, funcMap := range funcMaps {
		t.Funcs(funcMap)
	}

	t, err := t.Parse(templateStr)
	if err != nil {
		log.Warnf("Template parse file failed: %s.", err)
		return "", err
	}

	var buff bytes.Buffer
	if err = t.Execute(&buff, variable); err != nil {
		log.Debugf("Template execution failed: %s.", err)
		return "", err
	}
	return buff.String(), nil
}

func RenderForFile(name, tplFileName, dstFileName string, variable any) error {
	log.Debugf("Render filePath: %s.", tplFileName)
	log.Debugf("Render config: %v.", variable)
	log.Debugf("Render output: %s.", dstFileName)

	textBytes, err := os.ReadFile(tplFileName)
	if err != nil {
		return err
	}
	textStr := string(textBytes)
	renderedStr, err := Render(name, string(textStr), variable)
	if err != nil {
		log.Debugf("render %s failed: %s", name, err)
		return err
	}
	return os.WriteFile(dstFileName, []byte(renderedStr), 0644)
}

func DefaultRender(templateName string, vars any, funcMaps ...template.FuncMap) RenderFunc {
	return func(src []byte) (string, error) {
		return Render(templateName, string(src), vars, funcMaps...)
	}
}

func (r *rendererWithGetter) SetDefaultRender(templateName string, vars any, funcMaps ...template.FuncMap) *rendererWithRender {
	return r.SetRender(DefaultRender(templateName, vars, funcMaps...))
}
