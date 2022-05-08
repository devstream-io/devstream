package statemanager

import (
	"bytes"
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz/concurrentmap"
)

// State is the single component's state.
type State struct {
	Name       string
	InstanceID string
	DependsOn  []string
	Options    map[string]interface{}
	Resource   map[string]interface{}
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

// Note: Please use the StateKeyGenerateFunc function to generate StateKey instance.
type StateKey string

func StateKeyGenerateFunc(t *configloader.Tool) StateKey {
	return StateKey(fmt.Sprintf("%s_%s", t.Name, t.InstanceID))
}

func GenerateStateKeyByToolNameAndPluginKind(toolName string, instanceID string) StateKey {
	return StateKey(fmt.Sprintf("%s_%s", toolName, instanceID))
}
