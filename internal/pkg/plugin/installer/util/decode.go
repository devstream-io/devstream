package util

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

func DecodePlugin(rawOptions configmanager.RawOptions, pluginData any) error {
	// 1. create a new decode with pluginDecoder config
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: pluginDecoder,
		Result:     pluginData,
	})
	if err != nil {
		return fmt.Errorf("create plugin decoder failed: %w", err)
	}
	// 2. decode rawOptions to structData
	if err := decoder.Decode(rawOptions); err != nil {
		return fmt.Errorf("decode plugin option failed: %w", err)
	}
	return nil
}

func pluginDecoder(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	switch t {
	// set git.RepoInfo default value
	case reflect.TypeOf(&git.RepoInfo{}):
		repoData := new(git.RepoInfo)
		if err := mapstructure.Decode(data, repoData); err != nil {
			return nil, err
		}
		if err := repoData.SetDefault(); err != nil {
			return nil, err
		}
		return mapz.DecodeStructToMap(repoData)
	}
	return data, nil
}
