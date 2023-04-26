package patch

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/devstream-io/devstream/internal/log"
)

// Patch calls the patch command to apply a diff file to an original
func Patch(workDir, patchFile string) error {
	log.Infof("Patching file: %s", patchFile)

	// Check if the patch command exists and is executable
	err := checkPatchCommand()
	if err != nil {
		return fmt.Errorf("patch command check failed: %w", err)
	}

	// Use the patch tool to apply the patch
	cmd := exec.Command("patch", "-i", patchFile, "-t", "-p0")
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("patch command failed: %w\nOutput: %s", err, string(output))
	}

	log.Infof("Successfully patched the file")
	return nil
}

// checkPatchCommand checks if the patch command exists and is executable
func checkPatchCommand() error {
	// Check if the patch command exists
	path, err := exec.LookPath("patch")
	if err != nil {
		return fmt.Errorf("patch command not found: %w", err)
	}

	// Check if the patch command is executable
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat patch command: %w", err)
	}

	if fileInfo.Mode()&0111 == 0 {
		return fmt.Errorf("patch command is not executable")
	}

	return nil
}
