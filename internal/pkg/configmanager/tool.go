package configmanager

import (
	"fmt"
	"runtime"
	"strings"

	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

var (
	GOOS   string = runtime.GOOS
	GOARCH string = runtime.GOARCH
)

type RawOptions map[string]any

// Tool is the struct for one section of the DevStream tool file (part of the config.)
type Tool struct {
	Name string `yaml:"name" validate:"required"`
	// RFC 1123 - DNS Subdomain Names style
	// contain no more than 253 characters
	// contain only lowercase alphanumeric characters, '-' or '.'
	// start with an alphanumeric character
	// end with an alphanumeric character
	InstanceID string     `yaml:"instanceID" validate:"required,dns1123subdomain"`
	DependsOn  []string   `yaml:"dependsOn"`
	Options    RawOptions `yaml:"options"`
}

func newTool(name, instanceID string, options RawOptions) *Tool {
	if options == nil {
		options = make(RawOptions)
	}

	return &Tool{
		Name:       name,
		InstanceID: instanceID,
		DependsOn:  []string{},
		Options:    options,
	}
}

func (t *Tool) String() string {
	bs, err := yaml.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

type Tools []*Tool

func (tools Tools) validateAll() error {
	var errs []error
	errs = append(errs, tools.validate()...)
	errs = append(errs, tools.validateDependsOnConfig()...)
	return multierr.Combine(errs...)
}

func (tools Tools) validate() (errs []error) {
	for _, tool := range tools {
		errs = append(errs, tool.validate()...)
	}
	errs = append(errs, tools.duplicatedCheck()...)
	return
}

func (t *Tool) validate() []error {
	return validator.Struct(t)
}

func (t *Tool) DeepCopy() *Tool {
	var retTool = Tool{
		Name:       t.Name,
		InstanceID: t.InstanceID,
		DependsOn:  t.DependsOn,
		Options:    RawOptions{},
	}
	for k, v := range t.Options {
		retTool.Options[k] = v
	}
	return &retTool
}

func (t *Tool) KeyWithNameAndInstanceID() string {
	return fmt.Sprintf("%s.%s", t.Name, t.InstanceID)
}

// GetPluginName return plugin name without file extensions
func (t *Tool) GetPluginName() string {
	return fmt.Sprintf("%s-%s-%s_%s", t.Name, GOOS, GOARCH, version.Version)
}

func (t *Tool) GetPluginNameWithOSAndArch(os, arch string) string {
	return fmt.Sprintf("%s-%s-%s_%s", t.Name, os, arch, version.Version)
}

// GetPluginFileName creates the file name based on the tool's name and version
// If the plugin {github-actions 0.0.1}, the generated name will be "github-actions_0.0.1.so"
func (t *Tool) GetPluginFileName() string {
	return t.GetPluginName() + ".so"
}

func (t *Tool) GetPluginFileNameWithOSAndArch(os, arch string) string {
	return t.GetPluginNameWithOSAndArch(os, arch) + ".so"
}

func (t *Tool) GetPluginMD5FileName() string {
	return t.GetPluginName() + ".md5"
}
func (t *Tool) GetPluginMD5FileNameWithOSAndArch(os, arch string) string {
	return t.GetPluginNameWithOSAndArch(os, arch) + ".md5"
}

func (tools Tools) duplicatedCheck() (errs []error) {
	list := make(map[string]struct{})
	for _, t := range tools {
		key := t.KeyWithNameAndInstanceID()
		if _, ok := list[key]; ok {
			errs = append(errs, fmt.Errorf("tool or app <%s> is duplicated", key))
		}
		list[key] = struct{}{}
	}
	return errs
}

// validateDependsOnConfig is used to validate all tools' DependsOn config
func (tools Tools) validateDependsOnConfig() (retErrs []error) {
	retErrs = make([]error, 0)
	toolKeySet := make(map[string]struct{})

	// record all tools' key with name.instanceID
	for _, tool := range tools {
		toolKeySet[tool.KeyWithNameAndInstanceID()] = struct{}{}
	}

	validateOneTool := func(tool *Tool) (errs []error) {
		errs = make([]error, 0)
		if len(tool.DependsOn) == 0 {
			return
		}
		for _, d := range tool.DependsOn {
			if strings.TrimSpace(d) == "" {
				continue
			}

			if _, ok := toolKeySet[d]; !ok {
				errs = append(errs, fmt.Errorf("tool %s's DependsOn %s doesn't exist", tool.InstanceID, d))
			}
		}
		return
	}

	for _, t := range tools {
		retErrs = append(retErrs, validateOneTool(t)...)
	}

	return
}

func (tools Tools) updateToolDepends(dependTools Tools) {
	for _, tool := range tools {
		for _, dependTool := range dependTools {
			tool.DependsOn = append(tool.DependsOn, dependTool.KeyWithNameAndInstanceID())
		}
	}
}

func (tools Tools) renderInstanceIDtoOptions() {
	for _, t := range tools {
		if t.Options == nil {
			t.Options = make(RawOptions)
		}
		t.Options["instanceID"] = t.InstanceID
	}
}
