package status

import (
	"bytes"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/internal/pkg/statemanager"
)

type Output struct {
	InstanceID string                 `yaml:"InstanceID"`
	Plugin     string                 `yaml:"Plugin"`
	Drifted    bool                   `yaml:"Drifted"`
	Options    map[string]interface{} `yaml:"Options"`
	Status     *Status                `yaml:"Status"`
}

type Status struct {
	// if ResourceStatusInState == ResourceStatusFromRead,
	// then InlineStatus is set to equals ResourceStatusInState and ResourceStatusFromRead.
	InlineStatus           statemanager.ResourceStatus `yaml:",inline,omitempty"`
	ResourceStatusInState  statemanager.ResourceStatus `yaml:"statusInState,omitempty"`
	ResourceStatusFromRead statemanager.ResourceStatus `yaml:"statusNow,omitempty"`
}

// If the resource has drifted, status.ResourceStatusInState & status.ResourceStatusFromRead must NOT be nil and status.InlineStatus should be nil.
// If the resource hasn't drifted, status.ResourceStatusInState & status.ResourceStatusFromRead should be nil and status.InlineStatus must NOT be nil.
func NewOutput(instanceID, plugin string, options configmanager.RawOptions, status *Status) (*Output, error) {
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

func validateParams(instanceID, plugin string, options configmanager.RawOptions, status *Status) (bool, error) {
	if instanceID == "" || plugin == "" {
		return false, fmt.Errorf("instanceID or plugin cannot be nil")
	}
	if options == nil {
		return false, fmt.Errorf("options cannot be nil")
	}

	if status == nil {
		return false, fmt.Errorf("status cannot be nil")
	}
	if status.InlineStatus != nil && (status.ResourceStatusInState != nil || status.ResourceStatusFromRead != nil) {
		return false, fmt.Errorf("illegal status content")
	}
	if status.InlineStatus == nil && (status.ResourceStatusInState == nil || status.ResourceStatusFromRead == nil) {
		return false, fmt.Errorf("illegal status content")
	}

	return true, nil
}
