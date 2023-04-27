package commit

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/devstream-io/devstream/internal/log"
)

// Commit is used to execute git commit operations
func Commit(message string) error {
	// Check if the git command exists
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git command not found: %w", err)
	}

	cmd := exec.Command(gitPath, "commit", "-m", message)
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		return fmt.Errorf("git commit failed: %w\nOutput: %s", err, outputStr)
	}

	log.Infof("Successfully committed the file")
	return nil
}
