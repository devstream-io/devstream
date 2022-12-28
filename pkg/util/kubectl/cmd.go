package kubectl

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	Create string = "create"
	Apply  string = "apply"
	Delete string = "delete"
)

// KubeCreate runs "kubectl create -f filename"
func KubeCreate(filename string) error {
	return kubectlAction(Create, filename)
}

// KubeCreateFromIOReader generates a temp file from io.Reader and runs "kubectl create -f filename"
func KubeCreateFromIOReader(reader io.Reader) error {
	return ioToFileWrapper(reader, KubeCreate)
}

// KubeApply runs "kubectl apply -f filename"
func KubeApply(filename string) error {
	return kubectlAction(Apply, filename)
}

// KubeApplyFromIOReader generates a temp file from io.Reader and runs "kubectl apply -f filename"
func KubeApplyFromIOReader(reader io.Reader) error {
	return ioToFileWrapper(reader, KubeApply)
}

// KubeDelete runs "kubectl delete -f filename"
func KubeDelete(filename string) error {
	return kubectlAction(Delete, filename)
}

// KubeDeleteFromIOReader generates a temp file from io.Reader and runs "kubectl delete -f filename"
func KubeDeleteFromIOReader(reader io.Reader) error {
	return ioToFileWrapper(reader, KubeDelete)
}

func kubectlAction(action string, filename string) error {
	cmd := exec.Command("kubectl", action, "-f", filename)
	cOut, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to exec: < %s >.\nExec logs: < %s >. Got error: %w", cmd.String(), string(cOut), err)
	}
	log.Info(strings.TrimSuffix(string(cOut), "\n"))
	return nil
}

const defaultTempName = "kubectl_temp"

func ioToFileWrapper(reader io.Reader, f func(filename string) error) error {
	tempFile, err := os.CreateTemp("", defaultTempName)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, reader)
	if err != nil {
		return err
	}
	return f(tempFile.Name())
}
