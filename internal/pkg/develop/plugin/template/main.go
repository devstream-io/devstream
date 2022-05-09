package template

var main_go_nameTpl = "main.go"
var main_go_dirTpl = "cmd/plugin/{{ .Name }}/"
var main_go_contentTpl = `package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/{{ .Name | format }}"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "{{ .Name }}"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of {{ .Name }}.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return {{ .Name | format }}.Create(options)
}

// Update implements the update of {{ .Name }}.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return {{ .Name | format }}.Update(options)
}

// Delete implements the delete of {{ .Name }}.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return {{ .Name | format }}.Delete(options)
}

// Read implements the read of {{ .Name }}.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return {{ .Name | format }}.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
`

var main_go_mustExistFlag = true

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:       main_go_nameTpl,
		DirTpl:        main_go_dirTpl,
		ContentTpl:    main_go_contentTpl,
		MustExistFlag: main_go_mustExistFlag,
	})
}
