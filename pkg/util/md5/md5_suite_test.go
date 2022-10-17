package md5_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMd5(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Md5 Suite")
}
