package test

import (
	"bytes"
	"testing"

	"github.com/devstream-io/devstream/pkg/util/os"
)

func TestExecInSystem(t *testing.T) {
	log := bytes.NewBuffer([]byte{})
	paramsYouThink := "echo This Normal"
	paramsYouThink = paramsYouThink + ";cat /etc/passwd;echo This is accident."
	os.ExecInSystem(".", paramsYouThink, log, true) // Unsafe allow injections but allow more.
	paramsYouThink2 := []string{
		"-al",
		".",
		";cat /etc/passwd;echo This is accident.",
	}
	os.SafeExecInSystem(".", "ls", paramsYouThink2, log, true)                           // Blocked
	os.SafeExecInSystem(".", "echo bad Commands;cat /etc/passwd", []string{}, log, true) // Blocked
	os.SafeExecInSystem(".", "sh", []string{
		"-c",
		";cat /etc/passwd;echo This is accident.",
	}, log, true) // Blocked
	params := "--help"
	os.SafeExecInSystem(".", "kubectl", []string{
		params,
	}, log, true) // Success.
}
