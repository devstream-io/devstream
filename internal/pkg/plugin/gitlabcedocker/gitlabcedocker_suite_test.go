package gitlabcedocker_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGitlabcedocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gitlabcedocker Suite")
}
