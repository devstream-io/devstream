package jenkins

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func howToGetPasswdOfAdmin(opts *Options) {
	log.Info("Here is how to get the password of the admin user:")
	method := fmt.Sprintf("kubectl exec --namespace jenkins -it svc/%s-jenkins -c jenkins "+
		"-- /bin/cat /run/secrets/additional/chart-admin-password && echo", opts.Chart.ReleaseName)
	log.Info(method)
}

func showJenkinsUrl(opts *Options) {
	commands := []string{
		`jsonpath="{.spec.ports[0].nodePort}"`,
		fmt.Sprintf(`NODE_PORT=$(kubectl get -n jenkins -o jsonpath=$jsonpath services %s-jenkins)`, opts.Chart.ReleaseName),
		`jsonpath="{.items[0].status.addresses[0].address}"`,
		`NODE_IP=$(kubectl get nodes -n jenkins -o jsonpath=$jsonpath)`,
		`echo http://$NODE_IP:$NODE_PORT/login`,
	}

	cmd := exec.Command("sh", "-c", strings.Join(commands, " && "))

	output, err := cmd.Output()
	if err != nil {
		log.Errorf("Failed to get jenkins url: %v", err)
	}

	log.Info("Jenkins url:", string(output))

}
