package test

import (
	"bytes"
	"testing"

	"github.com/devstream-io/devstream/pkg/util/os"
)

func TestExecInSystem(t *testing.T) {
	log := bytes.NewBuffer([]byte{})
	// normal inject
	paramsYouThink := "echo This Normal"
	paramsYouThink = paramsYouThink + ";cat /etc/passwd;echo This is accident."
	os.ExecInSystem(".", paramsYouThink, log, true) // Unsafe allow injections but allow complex commands when hardcode.

	// inject to params
	paramsYouThink2 := []string{
		"-al",
		".",
		";cat /etc/passwd;echo This is accident.",
	}
	os.SafeExecInSystem(".", "ls", paramsYouThink2, log, true)                           // Blocked

	// inject to cmdName
	os.SafeExecInSystem(".", "echo bad Commands;cat /etc/passwd", []string{}, log, true) // Blocked

	// Try to introduce a shell command
	os.SafeExecInSystem(".", "sh", []string{
		"-c",
		";cat /etc/passwd;echo This is accident.",
	}, log, true) // Blocked
	params := "--help"

	// Try to test normal func.
	os.SafeExecInSystem(".", "kubectl", []string{
		params,
	}, log, true) // Success.
}
