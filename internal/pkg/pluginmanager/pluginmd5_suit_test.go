package pluginmanager

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPluginmanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Util Pluginmanager Suite")
}
