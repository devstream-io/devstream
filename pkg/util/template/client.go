package template

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// RenderAPI is a interface with render methods
type RenderAPI interface {
	Render(inputStr string, variables any) (string, error)
}

type (
	renderClient struct {
		Getter     getFunc
		Processors []processFunc
		// render config
		Option *TemplateOption
	}
	// TemplateOption is used to config template option
	TemplateOption struct {
		Name               string
		IgnoreMissKeyError bool
		FuncMap            template.FuncMap
	}
)

// NewRenderClient create template render client with template option, getter and processors
func NewRenderClient(clientOption *TemplateOption, getter getFunc, processors ...processFunc) RenderAPI {
	return &renderClient{
		Option:     clientOption,
		Getter:     getter,
		Processors: processors,
	}
}

// Render gets the content, process the content, render and returns the result string
func (c *renderClient) Render(inputStr string, variables any) (string, error) {
	// 1. get content
	content, err := c.Getter(inputStr)
	if err != nil {
		return "", err
	}
	// 2. process content to update content data
	for _, processFunc := range c.Processors {
		content = processFunc(content)
	}
	// 3. render content
	return c.renderContent(string(content), variables)
}

// renderContent render data with variable
func (c *renderClient) renderContent(templateStr string, variable any) (string, error) {
	t, err := c.newTemplateClient().Parse(templateStr)
	if err != nil {
		log.Warnf("Template parse file failed, template: %s, err: %s.", templateStr, err)
		return "", fmt.Errorf("parse %w", err)
	}

	var buff bytes.Buffer
	if err = t.Execute(&buff, variable); err != nil {
		log.Warnf("Template execution failed: %s.", err)
		return "", fmt.Errorf("render %w", err)
	}
	return buff.String(), nil
}

// newTemplateClient template client from  renderClient.Option
func (c *renderClient) newTemplateClient() *template.Template {
	// config options default value
	if c.Option == nil {
		c.Option = &TemplateOption{
			Name: "default_template",
		}
	}
	t := template.New(c.Option.Name).Delims("[[", "]]")
	if !c.Option.IgnoreMissKeyError {
		t = t.Option("missingkey=error")
	}
	// use sprig functions such as "env"
	t.Funcs(sprig.TxtFuncMap())
	if c.Option.FuncMap != nil {
		t.Funcs(c.Option.FuncMap)
	}
	return t
}
