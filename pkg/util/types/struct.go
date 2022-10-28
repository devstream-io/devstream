package types

import (
	"reflect"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"

	"github.com/mitchellh/mapstructure"
)

// FillStructDefaultValue will fill structData filed value by defaultStructData field value if they are empty
func FillStructDefaultValue(structData, defaultStructData any) {
	if defaultStructData == nil {
		return
	}
	defaultFalse := reflect.ValueOf(Bool(false))
	dataStruct := reflect.ValueOf(structData)
	structIterator := reflect.Indirect(dataStruct)
	defaultStruct := reflect.Indirect(reflect.ValueOf(defaultStructData))
	for i := 0; i < structIterator.NumField(); i++ {
		field := structIterator.Type().Field(i)
		fieldName := field.Name
		fieldVal := structIterator.Field(i).Interface()
		fieldType := structIterator.Field(i).Type()
		zeroValue := reflect.Zero(fieldType).Interface()

		if reflect.DeepEqual(fieldVal, zeroValue) {
			defaultValue := defaultStruct.FieldByName(fieldName)
			if !reflect.DeepEqual(defaultValue.Interface(), zeroValue) {
				dataStruct.Elem().Field(i).Set(defaultValue)
			} else if defaultValue.Kind() == reflect.Ptr {
				//TODO(steinliber): add more pointer judgement
				dataStruct.Elem().Field(i).Set(defaultFalse)
			}
		}
	}
}

// EncodeStruct will get structData and encode this data to map
func EncodeStruct(structData any) (configmanager.RawOptions, error) {
	var options configmanager.RawOptions
	if err := mapstructure.Decode(structData, &options); err != nil {
		return nil, err
	}
	return options, nil
}
