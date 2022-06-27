package gitlabcedocker

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type sshDockerOperator struct{}

func (op *sshDockerOperator) IfImageExists(imageNameWithTag string) bool {
	// eg. docker image ls gitlab/gitlab-ce:rc
	cmdString := fmt.Sprintf("docker image ls %v", imageNameWithTag)
	// output: eg.
	// REPOSITORY         TAG       IMAGE ID       CREATED      SIZE
	// gitlab/gitlab-ce   rc        a8543d702e39   4 days ago   2.49GB
	outputBuffer := &bytes.Buffer{}
	err := ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return false
	}
	// eg. gitlab/gitlab-ce
	imageNameWithoutTag := extractImageName(imageNameWithTag)

	return strings.Contains(outputBuffer.String(), imageNameWithoutTag)

}

func extractImageName(imageNameWithTag string) string {
	// the imageNameWithTag is in the format of "registry/image:tag"
	// we only want to return the image name "registry/image"
	return strings.Split(imageNameWithTag, ":")[0]
}

func (op *sshDockerOperator) PullImage(imageName string) error {
	err := ExecInSystemWithParams(".", []string{"docker", "pull", imageName}, nil, true)

	return err
}

func (op *sshDockerOperator) RemoveImage(imageName string) error {
	log.Infof("Removing image %v ...", imageName)

	cmdString := fmt.Sprintf("docker rmi %s", imageName)
	err := ExecInSystem(".", cmdString, nil, true)

	return err
}

func (op *sshDockerOperator) IfContainerExists(containerName string) bool {
	cmdString := fmt.Sprintf("docker inspect %s", containerName)
	outputBuffer := &bytes.Buffer{}
	err := ExecInSystem(".", cmdString, outputBuffer, false)
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
	err := ExecInSystem(".", cmdStringOneline, nil, true)
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
		options.HTTPPort, options.SSHPort, gitlabContainerName, options.GitLabHome, getImageNameWithTag(options))
	return cmdString
}

func (op *sshDockerOperator) StopContainer(containerName string) error {
	log.Infof("Stopping container %v ...", containerName)

	cmdString := fmt.Sprintf("docker stop %s", containerName)
	err := ExecInSystem(".", cmdString, nil, true)

	return err
}

func (op *sshDockerOperator) RemoveContainer(containerName string) error {
	log.Infof("Removing container %v ...", containerName)

	cmdString := fmt.Sprintf("docker rm %s", containerName)
	err := ExecInSystem(".", cmdString, nil, true)

	return err
}

func (op *sshDockerOperator) ListContainerMounts(containerName string) ([]string, error) {
	cmdString := fmt.Sprintf(`docker inspect --format='{{range .Mounts}}{{.Source}}{{"\n"}}{{end}}' %s`, containerName)
	outputBuffer := &bytes.Buffer{}

	err := ExecInSystem(".", cmdString, outputBuffer, false)
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

func (op *sshDockerOperator) GetContainerHostname(container string) (string, error) {
	cmdString := fmt.Sprintf("docker inspect --format='{{.Config.Hostname}}' %s", container)
	outputBuffer := &bytes.Buffer{}

	err := ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimSpace(outputBuffer.String())), nil

}

func (op *sshDockerOperator) GetContainerPortBinding(container, containerPort, protocol string) (hostPort string, err error) {

	// get container port binding map
	// the result is like:
	// 22/tcp->8122
	// 443/tcp->8443
	// 80/tcp->8180
	format := "'{{range $p,$conf := .NetworkSettings.Ports}}{{$p}}->{{(index $conf 0).HostPort}}{{println}}{{end}}'"
	cmdString := fmt.Sprintf("docker inspect --format=%s %s", format, container)
	outputBuffer := &bytes.Buffer{}
	err = ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return "", err
	}
	portBindings := strings.Split(strings.TrimSpace(outputBuffer.String()), "\n")
	log.Debugf("Container %v port bindings: %v", container, portBindings)

	// transfer port bindings to map
	portBindingsMap := make(map[string]string)
	for _, portBinding := range portBindings {
		portBindingParts := strings.Split(portBinding, "->")
		if len(portBindingParts) != 2 {
			return "", fmt.Errorf("Invalid port binding: %v", portBinding)
		}
		portBindingsMap[portBindingParts[0]] = portBindingParts[1]
	}

	portKey := fmt.Sprintf("%s/%s", containerPort, protocol)
	hostPort, ok := portBindingsMap[portKey]
	if !ok {
		return "", fmt.Errorf("No port binding for %v", portKey)
	}

	return hostPort, nil

}
