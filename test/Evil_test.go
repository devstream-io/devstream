package test

import (
	"bytes"
	"github.com/devstream-io/devstream/pkg/util/os"
	"testing"
)

func TestExecInSystem(t *testing.T) {
	paramsYouThink := []string{
		"whoami",
		"",
		";echo YOU GOT HACKED;",
	}
	log := bytes.NewBuffer([]byte{})
	os.ExecInSystem(".", paramsYouThink, log, true)
	os.SafeExecInSystem(".", "whoami", paramsYouThink, log, true)
}
