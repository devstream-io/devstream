package main

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/openldap"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// NAME is the name of this DevStream plugin.
const NAME = "openldap"

// Plugin is the type used by DevStream core. It's a string.
type Plugin string

// Create implements the create of OpenLDAP.
func (p Plugin) Create(options map[string]interface{}) (map[string]interface{}, error) {
	return openldap.Create(options)
}

// Update implements the update of OpenLDAP.
func (p Plugin) Update(options map[string]interface{}) (map[string]interface{}, error) {
	return openldap.Update(options)
}

// Delete implements the delete of OpenLDAP.
func (p Plugin) Delete(options map[string]interface{}) (bool, error) {
	return openldap.Delete(options)
}

// Read implements the read of OpenLDAP.
func (p Plugin) Read(options map[string]interface{}) (map[string]interface{}, error) {
	return openldap.Read(options)
}

// DevStreamPlugin is the exported variable used by the DevStream core.
var DevStreamPlugin Plugin

func main() {
	log.Infof("%T: %s is a plugin for DevStream. Use it with DevStream.\n", NAME, DevStreamPlugin)
}
