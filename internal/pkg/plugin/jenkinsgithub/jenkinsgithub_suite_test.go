package jenkinsgithub

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJenkinsgithub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jenkinsgithub Suite")
}
