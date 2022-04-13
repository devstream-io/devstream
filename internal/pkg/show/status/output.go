package status

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Output struct {
	Name    string                 `yaml:"Name"`
	Plugin  string                 `yaml:"Plugin"`
	Drifted bool                   `yaml:"Drifted"`
	Options map[string]interface{} `yaml:"Options"`
	Status  *Status                `yaml:"Status"`
}

type Status struct {
	InnerStatus map[string]interface{} `yaml:",inline,omitempty"`
	State       map[string]interface{} `yaml:"State,omitempty"`
	Resource    map[string]interface{} `yaml:"Resource,omitempty"`
}

// If the resource has drifted, status.State & status.Resource must NOT be nil and status.InnerStatus should be nil.
// If the resource hasn't drifted, status.State & status.Resource should be nil and status.InnerStatus must NOT be nil.
func NewOutput(name, plugin string, options map[string]interface{}, status *Status) (*Output, error) {
	// params validation
	if name == "" || plugin == "" {
		return nil, fmt.Errorf("name or plugin cannot be nil")
	}
	if options == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	if status == nil {
		return nil, fmt.Errorf("status cannot be nil")
	}
	if status.InnerStatus != nil && (status.State != nil || status.Resource != nil) {
		return nil, fmt.Errorf("illegal status content")
	}
	if status.InnerStatus == nil && (status.State == nil || status.Resource == nil) {
		return nil, fmt.Errorf("illegal status content")
	}

	// fill the output object
	output := &Output{
		Name:    name,
		Plugin:  plugin,
		Options: options,
		Status:  status,
	}

	if status.InnerStatus != nil {
		output.Drifted = false
	}

	return output, nil
}

func (o *Output) Print() error {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	err := encoder.Encode(o)
	if err != nil {
		return err
	}

	fmt.Println(buf.String())
	return nil
}
