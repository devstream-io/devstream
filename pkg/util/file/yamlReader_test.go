package file

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
)

func TestYamlReaderFile(t *testing.T) {
	content, err := ReadYamls("config-gitops.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(content))
}

func TestYamlReaderDir(t *testing.T) {
	content, err := ReadYamls("tmp")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(content))
}

var _ = Describe("ReadYamls func", func() {

})
