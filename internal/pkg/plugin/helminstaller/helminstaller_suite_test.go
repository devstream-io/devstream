package helminstaller_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHelmInstaller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "helminstaller suite")
}
