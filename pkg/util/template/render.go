package template

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func Render(name, templateStr string, variable any) (string, error) {
	t, err := template.New(name).Delims("[[", "]]").Parse(templateStr)
	if err != nil {
		log.Debugf("Template parse file failed: %s.", err)
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

	textBytes, err := ioutil.ReadFile(tplFileName)
	if err != nil {
		return err
	}
	textStr := string(textBytes)
	renderedStr, err := Render(name, string(textStr), variable)
	if err != nil {
		log.Debugf("render %s failed: %s", name, err)
		return err
	}
	return ioutil.WriteFile(dstFileName, []byte(renderedStr), 0666)
}
