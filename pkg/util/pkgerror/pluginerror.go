package pkgerror

import (
	"fmt"
	"strings"
)

type ErrorMessage string

// PluginError is plugin's error wrapper
type PluginError struct {
	// which plugin the error raises
	PluginName string
	// error message
	Message string
	// error extra module info
	Module string
}

func (e *PluginError) Error() string {
	moduleMsg := ""
	if e.Module != "" {
		moduleMsg = fmt.Sprintf("(%s)", e.Module)
	}
	return fmt.Sprintf("[%s]%s: %s", e.PluginName, moduleMsg, e.Message)
}

func CheckErrorMatchByMessage(targetError error, errCheckStrings ...ErrorMessage) bool {
	for _, checkErr := range errCheckStrings {
		if strings.Contains(targetError.Error(), string(checkErr)) {
			return true
		}
	}
	return false
}

func NewErrorFromPlugin(pluginName, moduleName string, err error) *PluginError {
	return &PluginError{
		PluginName: pluginName,
		Message:    err.Error(),
		Module:     moduleName,
	}
}
