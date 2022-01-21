package statemanager

import (
	"bytes"

	"github.com/merico-dev/stream/internal/pkg/log"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/container/concurrentmap"
)

type StatesMap concurrentmap.ConcurrentMap

func NewStatesMap() *StatesMap {
	return (*StatesMap)(concurrentmap.NewConcurrentMap("", new(State)))
}

func (s StatesMap) DeepCopy() *StatesMap {
	newStatesMap := NewStatesMap()
	s.Range(func(key, value interface{}) bool {
		newStatesMap.Store(key, value)
		return true
	})
	return newStatesMap
}

func (s StatesMap) Format() []byte {
	tmpMap := make(map[string]*State)
	s.Range(func(key, value interface{}) bool {
		tmpMap[key.(string)] = value.(*State)
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
