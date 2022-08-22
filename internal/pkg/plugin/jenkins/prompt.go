package jenkins

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func buildPasswdOfAdminCommand(opts jenkinsOptions) string {
	jenkinsFullName := opts.getJenkinsFullName()

	return fmt.Sprintf("kubectl exec --namespace %s -it svc/%s -c jenkins "+
		"-- /bin/cat /run/secrets/additional/chart-admin-password && echo", opts.Chart.Namespace, jenkinsFullName)
}

func howToGetPasswdOfAdmin(options plugininstaller.RawOptions) error {
	opts, err := newOptions(options)
	if err != nil {
		return err
	}

	log.Info("Here is how to get the password of the admin user:")
	command := buildPasswdOfAdminCommand(opts)
	log.Info(command)

	return nil
}

func getPasswdOfAdmin(options plugininstaller.RawOptions) (string, error) {
	opts, err := newOptions(options)
	if err != nil {
		return "", err
	}

	commandString := buildPasswdOfAdminCommand(opts)
	command := exec.Command("sh", "-c", commandString)

	password, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get password of admin user: %v", err)
	}

	return strings.TrimSpace(string(password)), nil
}
