package jenkins

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func buildPasswdOfAdminCommand(opts jenkinsOptions) string {
	method := fmt.Sprintf("kubectl exec --namespace jenkins -it svc/%s-jenkins -c jenkins "+
		"-- /bin/cat /run/secrets/additional/chart-admin-password && echo", opts.Chart.ReleaseName)

	return method
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

// getJenkinsURL returns the jenkins url of the jenkins, format: hostname:port
func getJenkinsURL(options plugininstaller.RawOptions) (string, error) {
	opts, err := newOptions(options)
	if err != nil {
		return "", err
	}

	commands := []string{
		`jsonpath="{.spec.ports[0].nodePort}"`,
		fmt.Sprintf(`NODE_PORT=$(kubectl get -n jenkins -o jsonpath=$jsonpath services %s-jenkins)`, opts.Chart.ReleaseName),
		`jsonpath="{.items[0].status.addresses[0].address}"`,
		`NODE_IP=$(kubectl get nodes -n jenkins -o jsonpath=$jsonpath)`,
		`echo $NODE_IP:$NODE_PORT`,
	}

	cmd := exec.Command("sh", "-c", strings.Join(commands, " && "))

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to get jenkins url: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// getJenkinsURL behaves like getJenkinsURL, but the hostname will be always 127.0.0.1
func getJenkinsURLForTestEnv(options plugininstaller.RawOptions) (string, error) {
	opts, err := newOptions(options)
	if err != nil {
		return "", err
	}

	commands := []string{
		`jsonpath="{.spec.ports[0].nodePort}"`,
		fmt.Sprintf(`NODE_PORT=$(kubectl get -n jenkins -o jsonpath=$jsonpath services %s-jenkins)`, opts.Chart.ReleaseName),
		`jsonpath="{.items[0].status.addresses[0].address}"`,
		`NODE_IP=127.0.0.1`,
		`echo $NODE_IP:$NODE_PORT`,
	}

	cmd := exec.Command("sh", "-c", strings.Join(commands, " && "))

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("Failed to get jenkins url: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func showJenkinsUrl(options plugininstaller.RawOptions) error {
	opts, err := newOptions(options)
	if err != nil {
		return err
	}

	// prod env: just print Jenkins url
	if !opts.TestEnv {
		url, err := getJenkinsURL(options)
		if err != nil {
			log.Error(err)
			return err
		}

		log.Infof("Jenkins url: http://%s/login", url)
	}

	// test env: print Jenkins url in host machine and Jenkins url in K8s cluster
	if opts.TestEnv {
		log.Info("You are in test env. Here are the Jenkins url in host machine and Jenkins url in K8s cluster.")

		urlForTestEnv, err := getJenkinsURLForTestEnv(options)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Infof("Jenkins url in host machine: http://%s/login", urlForTestEnv)

		urlInK8s, err := getJenkinsURL(options)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Info("Jenkins url in K8s:", fmt.Sprintf("http://%s/login", urlInK8s))

	}

	return nil
}
