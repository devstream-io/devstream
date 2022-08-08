package dockersh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// ShellOperator is an implementation of /pkg/util/docker.Operator interface by using shell commands
type ShellOperator struct{}

func (op *ShellOperator) ImageIfExist(imageNameWithTag string) bool {
	// eg. docker image ls gitlab/gitlab-ce:rc -q
	// output: image id (if exist)
	cmdString := fmt.Sprintf("docker image ls %v -q", imageNameWithTag)
	outputBuffer := &bytes.Buffer{}
	err := ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return false
	}

	return strings.TrimSpace(outputBuffer.String()) != ""

}

func (op *ShellOperator) ImagePull(imageNameWithTag string) error {
	err := ExecInSystemWithParams(".", []string{"docker", "pull", imageNameWithTag}, nil, true)

	return err
}

func (op *ShellOperator) ImageRemove(imageNameWithTag string) error {
	log.Infof("Removing image %v ...", imageNameWithTag)

	cmdString := fmt.Sprintf("docker rmi %s", imageNameWithTag)
	err := ExecInSystem(".", cmdString, nil, true)

	return err
}

func (op *ShellOperator) ContainerIfExist(containerName string) bool {
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

func (op *ShellOperator) ContainerIfRunning(containerName string) bool {
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

func (op *ShellOperator) ContainerRun(opts *docker.RunOptions) error {
	// build the command
	cmdString, err := BuildContainerRunCommand(opts)
	if err != nil {
		return err
	}
	log.Debugf("Docker run command: %s", cmdString)

	// run the command
	err = ExecInSystem(".", cmdString, nil, true)
	if err != nil {
		return fmt.Errorf("docker run failed: %v", err)
	}

	return nil
}

// BuildContainerRunCommand builds the docker run command string from the given options and additional params
func BuildContainerRunCommand(opts *docker.RunOptions) (string, error) {
	if err := opts.Validate(); err != nil {
		return "", err
	}

	cmdBuilder := strings.Builder{}
	cmdBuilder.WriteString("docker run --detach ")
	if opts.Hostname != "" {
		cmdBuilder.WriteString(fmt.Sprintf("--hostname %s ", opts.Hostname))
	}
	for _, publish := range opts.PortPublishes {
		cmdBuilder.WriteString(fmt.Sprintf("--publish %d:%d ", publish.HostPort, publish.ContainerPort))
	}
	cmdBuilder.WriteString(fmt.Sprintf("--name %s ", opts.ContainerName))
	if opts.RestartAlways {
		cmdBuilder.WriteString("--restart always ")
	}
	for _, volume := range opts.Volumes {
		cmdBuilder.WriteString(fmt.Sprintf("--volume %s:%s ", volume.HostPath, volume.ContainerPath))
	}
	for _, param := range opts.RunParams {
		cmdBuilder.WriteString(param + " ")
	}
	cmdBuilder.WriteString(docker.CombineImageNameAndTag(opts.ImageName, opts.ImageTag))

	return cmdBuilder.String(), nil
}

func (op *ShellOperator) ContainerStop(containerName string) error {
	log.Infof("Stopping container %v ...", containerName)

	cmdString := fmt.Sprintf("docker stop %s", containerName)
	err := ExecInSystem(".", cmdString, nil, true)

	return err
}

func (op *ShellOperator) ContainerRemove(containerName string) error {
	log.Infof("Removing container %v ...", containerName)

	cmdString := fmt.Sprintf("docker rm %s", containerName)
	err := ExecInSystem(".", cmdString, nil, true)

	return err
}

func (op *ShellOperator) ContainerListMounts(containerName string) (docker.Mounts, error) {
	cmdString := fmt.Sprintf(`docker inspect --format='{{json .Mounts}}' %s`, containerName)

	outputBuffer := &bytes.Buffer{}

	err := ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return nil, err
	}

	mounts := make([]docker.MountPoint, 0)
	err = json.Unmarshal([]byte(strings.TrimSpace(outputBuffer.String())), &mounts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal docker inspect output when list mounts: %v", err)
	}

	log.Debugf("Container %v mounts: %v", containerName, mounts)

	return mounts, nil
}

func (op *ShellOperator) ContainerGetHostname(container string) (string, error) {
	cmdString := fmt.Sprintf("docker inspect --format='{{.Config.Hostname}}' %s", container)
	outputBuffer := &bytes.Buffer{}

	err := ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimSpace(outputBuffer.String())), nil

}

func (op *ShellOperator) ContainerListPortPublishes(containerName string) ([]docker.PortPublish, error) {
	// get container port binding map
	// the result is like:
	// 22/tcp->8122
	// 443/tcp->8443
	// 80/tcp->8180
	format := "'{{range $p,$conf := .NetworkSettings.Ports}}{{$p}}->{{(index $conf 0).HostPort}}{{println}}{{end}}'"
	cmdString := fmt.Sprintf("docker inspect --format=%s %s", format, containerName)
	outputBuffer := &bytes.Buffer{}
	err := ExecInSystem(".", cmdString, outputBuffer, false)
	if err != nil {
		return nil, err
	}
	portBindings := strings.Split(strings.TrimSpace(outputBuffer.String()), "\n")
	log.Debugf("Container %v port bindings: %v", containerName, portBindings)

	publishes, err := buildPortPublishes(portBindings)
	if err != nil {
		return publishes, err
	}

	return publishes, nil
}

func buildPortPublishes(portBindings []string) (PortPublishes []docker.PortPublish, err error) {
	// 22/tcp->8122
	// 443/tcp->8443
	// 80/tcp->8180
	re := regexp.MustCompile(`^(\d+)/(tcp|udp)->(\d+)$`)

	for _, portBinding := range portBindings {
		match := re.FindStringSubmatch(portBinding)
		// match e.g. ["22/tcp->8122", "22", "tcp", "8122"]
		if len(match) != 4 {
			return nil, fmt.Errorf("invalid port binding: %v", portBinding)
		}

		hostPort, err := strconv.Atoi(match[3])
		if err != nil {
			return nil, fmt.Errorf("invalid port binding: %v", portBinding)
		}
		containerPort, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("invalid port binding: %v", portBinding)
		}

		portPublish := docker.PortPublish{
			ContainerPort: uint(containerPort),
			HostPort:      uint(hostPort),
		}
		PortPublishes = append(PortPublishes, portPublish)
	}

	return PortPublishes, nil
}

func (op *ShellOperator) ContainerGetPortBinding(container string, containerPort uint) (hostPort uint, err error) {
	portBindings, err := op.ContainerListPortPublishes(container)
	if err != nil {
		return 0, err
	}

	for _, portBinding := range portBindings {
		if portBinding.ContainerPort == containerPort {
			return portBinding.HostPort, nil
		}
	}

	return 0, fmt.Errorf("container %v does not have port binding for port %v", container, containerPort)
}
