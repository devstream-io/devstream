package util_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

func TestCommon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugin Installer Util Suite")
}

type mockStruct struct {
	Scm        *git.RepoInfo `mapstructure:"scm"`
	DeepStruct deepStruct    `mapstructure:"deepStruct"`
}

type deepStruct struct {
	DeepStr string `mapstructure:"deepStr" validate:"required"`
}
