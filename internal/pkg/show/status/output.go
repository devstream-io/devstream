package status

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Output struct {
	InstanceID string                 `yaml:"InstanceID"`
	Plugin     string                 `yaml:"Plugin"`
	Drifted    bool                   `yaml:"Drifted"`
	Options    map[string]interface{} `yaml:"Options"`
	Status     *Status                `yaml:"Status"`
}

type Status struct {
	InlineStatus map[string]interface{} `yaml:",inline,omitempty"`
	State        map[string]interface{} `yaml:"State,omitempty"`
	Resource     map[string]interface{} `yaml:"Resource,omitempty"`
}

// If the resource has drifted, status.State & status.Resource must NOT be nil and status.InlineStatus should be nil.
// If the resource hasn't drifted, status.State & status.Resource should be nil and status.InlineStatus must NOT be nil.
func NewOutput(instanceID, plugin string, options map[string]interface{}, status *Status) (*Output, error) {
	if ok, err := validateParams(instanceID, plugin, options, status); !ok {
		return nil, err
	}

	output := &Output{
		InstanceID: instanceID,
		Plugin:     plugin,
		Drifted:    false,
		Options:    options,
		Status:     status,
	}

	if status.InlineStatus == nil {
		output.Drifted = true
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

func validateParams(instanceID, plugin string, options map[string]interface{}, status *Status) (bool, error) {
	if instanceID == "" || plugin == "" {
		return false, fmt.Errorf("instanceID or plugin cannot be nil")
	}
	if options == nil {
		return false, fmt.Errorf("options cannot be nil")
	}

	if status == nil {
		return false, fmt.Errorf("status cannot be nil")
	}
	if status.InlineStatus != nil && (status.State != nil || status.Resource != nil) {
		return false, fmt.Errorf("illegal status content")
	}
	if status.InlineStatus == nil && (status.State == nil || status.Resource == nil) {
		return false, fmt.Errorf("illegal status content")
	}

	return true, nil
}
