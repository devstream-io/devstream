package jenkinspipelinekubernetes

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJenkinspipelinekubernetes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jenkinspipelinekubernetes Suite")
}
