package statemanager

import (
	"bytes"
	"fmt"
	"sort"

	"golang.org/x/exp/maps"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz/concurrentmap"
)

// We call what the plugin created a ResourceStatus, and which is stored as part of the state.
type ResourceStatus map[string]interface{}

type ResourceOutputs map[string]interface{}

// State is the single component's state.
type State struct {
	Name           string                   `yaml:"name"`
	InstanceID     string                   `yaml:"instanceID"`
	DependsOn      []string                 `yaml:"dependsOn"`
	Options        configmanager.RawOptions `yaml:"options"`
	ResourceStatus ResourceStatus           `yaml:"resourceStatus"`
}

func (rs ResourceStatus) SetOutputs(outputs ResourceOutputs) {
	(rs)["outputs"] = outputs
}

func (rs ResourceStatus) GetOutputs() ResourceOutputs {
	outputs, ok := (rs)["outputs"]
	if !ok {
		return nil
	}
	outputStatus, isStatus := outputs.(ResourceStatus)
	if !isStatus {
		return outputs.(ResourceOutputs)
	}
	// if outputs type is ResourceStatus, transfer this type to ResourceOutputs
	outputData := ResourceOutputs{}
	maps.Copy(outputData, outputStatus)
	return outputData
}

type StatesMap struct {
	*concurrentmap.ConcurrentMap
}

func NewStatesMap() StatesMap {
	return StatesMap{
		ConcurrentMap: concurrentmap.NewConcurrentMap(StateKey(""), State{}),
	}
}

func (s StatesMap) DeepCopy() StatesMap {
	newStatesMap := NewStatesMap()
	s.Range(func(key, value interface{}) bool {
		newStatesMap.Store(key, value)
		return true
	})
	return newStatesMap
}

func (s StatesMap) ToList() []State {
	var res []State
	s.Range(func(key, value interface{}) bool {
		res = append(res, value.(State))
		return true
	})

	sort.Slice(res, func(i, j int) bool {
		keyi := fmt.Sprintf("%s.%s", res[i].InstanceID, res[i].Name)
		keyj := fmt.Sprintf("%s.%s", res[j].InstanceID, res[j].Name)
		return keyi < keyj
	})

	return res
}

func (s StatesMap) Format() []byte {
	tmpMap := make(map[StateKey]State)
	s.Range(func(key, value interface{}) bool {
		tmpMap[key.(StateKey)] = value.(State)
		return true
	})

	if len(tmpMap) == 0 {
		return []byte{}
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err := encoder.Encode(&tmpMap)
	if err != nil {
		log.Error(err)
		return nil
	}

	return buf.Bytes()
}

// Note: Please use the GenerateStateKeyByToolNameAndInstanceID function to generate StateKey instance.
type StateKey string

func GenerateStateKeyByToolNameAndInstanceID(toolName string, instanceID string) StateKey {
	return StateKey(fmt.Sprintf("%s_%s", toolName, instanceID))
}
