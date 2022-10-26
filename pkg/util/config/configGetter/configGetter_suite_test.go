package configGetter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfigGetter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ConfigGetter Suite")
}
