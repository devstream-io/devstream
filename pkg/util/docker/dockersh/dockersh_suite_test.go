package dockersh

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDockersh(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dockersh Suite")
}
