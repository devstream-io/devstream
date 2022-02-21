package pluginengine_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPluginengine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pluginengine Suite")
}
