package template

var main_go_nameTpl = "main.go"
var main_go_dirTpl = "cmd/{{ .Name }}/"
var main_go_contentTpl = `package main

import (
	"github.com/merico-dev/stream/internal/pkg/plugin/{{ .Name }}"
	"github.com/merico-dev/stream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "{{ .Name }}"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of {{ .Name }}.
func (p Plugin) Create(params map[string]interface{}) (map[string]interface{}, error) {
	return {{ .Name }}.Create(params)
}

// Update implements the update of {{ .Name }}.
func (p Plugin) Update(params map[string]interface{}) (map[string]interface{}, error) {
	return {{ .Name }}.Update(params)
}

// Delete implements the delete of {{ .Name }}.
func (p Plugin) Delete(params map[string]interface{}) (bool, error) {
	return {{ .Name }}.Delete(params)
}

// Read implements the read of {{ .Name }}.
func (p Plugin) Read(params map[string]interface{}) (map[string]interface{}, error) {
	return {{ .Name }}.Read(params)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    main_go_nameTpl,
		DirTpl:     main_go_dirTpl,
		ContentTpl: main_go_contentTpl,
	})
}
