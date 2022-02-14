package kubectl

import (
	"os/exec"
	"strings"

	"github.com/merico-dev/stream/internal/pkg/log"
)

const APPLY string = "apply"
const DELETE string = "delete"

// KubeApply runs "kubectl apply -f filename"
func KubeApply(filename string) error {
	return kubectlAction(APPLY, filename)
}

// KubeDelete runs "kubectl delete -f filename"
func KubeDelete(filename string) error {
	return kubectlAction(DELETE, filename)
}

func kubectlAction(action string, filename string) error {
	cmd := exec.Command("kubectl", action, "-f", filename)
	cOut, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("failed to exec: < %s >", cmd.String())
		log.Errorf("exec logs: < %s >. got error: %s", string(cOut), err)
		return err
	}
	log.Info(strings.TrimSuffix(string(cOut), "\n"))
	return nil
}
