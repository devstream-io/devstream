package concurrentmap

import (
	"fmt"
	"reflect"
	"sync"
)

type ConcurrentMap struct {
	*sync.Map
	KeyType   reflect.Type
	ValueType reflect.Type
}

func NewConcurrentMap(keyType, valueType interface{}) *ConcurrentMap {
	return &ConcurrentMap{
		Map:       &sync.Map{},
		KeyType:   reflect.TypeOf(keyType),
		ValueType: reflect.TypeOf(valueType),
	}
}

func (cm *ConcurrentMap) Load(key interface{}) (value interface{}, ok bool) {
	if reflect.TypeOf(key) == cm.KeyType {
		return cm.Map.Load(key)
	}
	return
}

func (cm *ConcurrentMap) Store(key, value interface{}) {
	if reflect.TypeOf(key) == cm.KeyType && reflect.TypeOf(value) == cm.ValueType {
		cm.Map.Store(key, value)
		return
	}
	panic(fmt.Errorf("wrong key or value type: %v, %v", reflect.TypeOf(key), reflect.TypeOf(value)))
}

func (cm *ConcurrentMap) Delete(key interface{}) {
	if reflect.TypeOf(key) == cm.KeyType {
		cm.Map.Delete(key)
		return
	}
	panic(fmt.Errorf("wrong key type: %v", reflect.TypeOf(key)))
}

func (cm *ConcurrentMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	if reflect.TypeOf(key) == cm.KeyType && reflect.TypeOf(value) == cm.ValueType {
		actual, loaded = cm.Map.LoadOrStore(key, value)
		return
	}
	panic(fmt.Errorf("wrong key or value type: %v, %v", reflect.TypeOf(key), reflect.TypeOf(value)))
}

func (cm *ConcurrentMap) Range(f func(key, value interface{}) bool) {
	cm.Map.Range(f)
}
