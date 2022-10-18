package template

var mainGoNameTpl = "main.go"
var mainGoDirTpl = "cmd/plugin/[[ .Name ]]/"
var mainGoContentTpl = `package main

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/[[ .Name | format ]]"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "[[ .Name ]]"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of [[ .Name ]].
func (p Plugin) Create(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return [[ .Name | format ]].Create(options)
}

// Update implements the update of [[ .Name ]].
func (p Plugin) Update(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return [[ .Name | format ]].Update(options)
}

// Delete implements the delete of [[ .Name ]].
func (p Plugin) Delete(options configmanager.RawOptions) (bool, error) {
	return [[ .Name | format ]].Delete(options)
}

// Read implements the read of [[ .Name ]].
func (p Plugin) Read(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	return [[ .Name | format ]].Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", DevStreamPlugin, NAME)
}
`

var mainGoMustExistFlag = true

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:       mainGoNameTpl,
		DirTpl:        mainGoDirTpl,
		ContentTpl:    mainGoContentTpl,
		MustExistFlag: mainGoMustExistFlag,
	})
}
