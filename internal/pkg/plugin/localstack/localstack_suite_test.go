package localstack

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLocalStackPlugin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LocalStack Plugin Suite")
}
