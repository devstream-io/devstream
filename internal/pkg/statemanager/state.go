package statemanager

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/pkg/util/log"
	"github.com/merico-dev/stream/pkg/util/mapz/concurrentmap"
)

// State is the single component's state.
type State struct {
	Name     string
	Plugin   configloader.Plugin
	Options  map[string]interface{}
	Resource map[string]interface{}
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
	return StateKey(fmt.Sprintf("%s_%s", t.Name, t.Plugin.Kind))
}

func GenStateKey(refParam []string) StateKey {
	return StateKey(fmt.Sprintf("%s_%s", refParam[0], refParam[1]))
}
