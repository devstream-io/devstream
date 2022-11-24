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

func (ts Tools) validateAll() error {
	var errs []error
	errs = append(errs, ts.validate()...)
	errs = append(errs, ts.validateDependency()...)
	return multierr.Combine(errs...)
}

func (ts Tools) validate() (errs []error) {
	for _, tool := range ts {
		errs = append(errs, tool.validate()...)
	}
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
// If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.so"
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

func (ts Tools) validateDependency() []error {
	errors := make([]error, 0)

	// config "set" (map)
	toolMap := make(map[string]bool)
	// creating the set
	for _, tool := range ts {
		toolMap[tool.KeyWithNameAndInstanceID()] = true
	}

	for _, tool := range ts {
		// no dependency, pass
		if len(tool.DependsOn) == 0 {
			continue
		}

		// for each dependency
		for _, dependency := range tool.DependsOn {
			// skip empty string
			dependency = strings.TrimSpace(dependency)
			if dependency == "" {
				continue
			}

			// generate an error if the dependency isn't in the config set,
			if _, ok := toolMap[dependency]; !ok {
				errors = append(errors, fmt.Errorf("tool %s's dependency %s doesn't exist in the config", tool.InstanceID, dependency))
			}
		}
	}

	return errors
}
