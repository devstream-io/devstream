package gitlabcedocker

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
	osUtil "github.com/devstream-io/devstream/pkg/util/os"
)

type sshDockerOperator struct{}

func (op *sshDockerOperator) IfImageExists(imageName string) bool {
	cmdString := fmt.Sprintf("docker image ls %v", imageName)
	outputBuffer := &bytes.Buffer{}
	err := osUtil.ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return false
	}

	return strings.Contains(outputBuffer.String(), imageName)

}

func (op *sshDockerOperator) PullImage(imageName string) error {
	err := osUtil.ExecInSystemWithParams(".", []string{"docker", "pull", imageName}, nil, true)

	return err
}

func (op *sshDockerOperator) RemoveImage(imageName string) error {
	log.Infof("Removing image %v ...", imageName)

	cmdString := fmt.Sprintf("docker rmi %s", imageName)
	err := osUtil.ExecInSystem(".", cmdString, nil, true)

	return err
}

func (op *sshDockerOperator) IfContainerExists(containerName string) bool {
	cmdString := fmt.Sprintf("docker inspect %s", containerName)
	outputBuffer := &bytes.Buffer{}
	err := osUtil.ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return false
	}

	if strings.Contains("No such object", outputBuffer.String()) {
		return false
	}

	return true
}

func (op *sshDockerOperator) IfContainerRunning(containerName string) bool {
	command := exec.Command("docker", "inspect", "--format='{{.State.Status}}'", containerName)
	output, err := command.Output()
	if err != nil {
		return false
	}

	if strings.Contains(string(output), "running") {
		return true
	}

	return false
}

func (op *sshDockerOperator) RunContainer(options Options) error {
	cmdString := BuildDockerRunCommand(options)
	log.Debugf("Docker run command: %s", cmdString)
	cmdStringOneline := strings.Replace(cmdString, "\\\n", " ", -1)
	err := osUtil.ExecInSystem(".", cmdStringOneline, nil, true)
	if err != nil {
		return fmt.Errorf("docker run failed: %v", err)
	}

	return nil
}

func BuildDockerRunCommand(options Options) string {
	cmdTemplate := `
	docker run --detach \
	--hostname %s \
	--publish %d:443 --publish %d:80 --publish %d:22 \
	--name %s \
	--restart always \
	--volume %[6]s/config:/etc/gitlab \
	--volume %[6]s/logs:/var/log/gitlab \
	--volume %[6]s/data:/var/opt/gitlab \
	--shm-size 256m \
	%s
	`
	cmdString := fmt.Sprintf(cmdTemplate, options.Hostname, options.HTTPSPort,
		options.HTTPPort, options.SSHPort, gitlabContainerName, options.GitLabHome, gitlabImageName)
	return cmdString
}

func (op *sshDockerOperator) StopContainer(containerName string) error {
	log.Infof("Stopping container %v ...", containerName)

	cmdString := fmt.Sprintf("docker stop %s", containerName)
	err := osUtil.ExecInSystem(".", cmdString, nil, true)

	return err
}

func (op *sshDockerOperator) RemoveContainer(containerName string) error {
	log.Infof("Removing container %v ...", containerName)

	cmdString := fmt.Sprintf("docker rm %s", containerName)
	err := osUtil.ExecInSystem(".", cmdString, nil, true)

	return err
}

func (op *sshDockerOperator) ListContainerMounts(containerName string) ([]string, error) {
	cmdString := fmt.Sprintf(`docker inspect --format='{{range .Mounts}}{{.Source}}{{"\n"}}{{end}}' %s`, containerName)
	outputBuffer := &bytes.Buffer{}

	err := osUtil.ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return nil, err
	}

	volumes := strings.Split(strings.TrimSpace(outputBuffer.String()), "\n")

	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i] < volumes[j]
	})

	log.Debugf("Container %v volumes: %v", containerName, volumes)

	return volumes, nil
}
