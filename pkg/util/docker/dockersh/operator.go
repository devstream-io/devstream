package dockersh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// ShellOperator is an implementation of /pkg/util/docker.Operator interface by using shell commands
type ShellOperator struct{}

// TODO(aFlyBird0): maybe use one param "ImageNameWithTag" is not a good idea,
// because we have to extract the bare image name from the image name with tag
// we could use (a struct)/(a interface)/(two params) to represent the image name and tag
func (op *ShellOperator) ImageIfExist(imageNameWithTag string) bool {
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

func (op *ShellOperator) ImagePull(imageName string) error {
	err := ExecInSystemWithParams(".", []string{"docker", "pull", imageName}, nil, true)

	return err
}

func (op *ShellOperator) ImageRemove(imageName string) error {
	log.Infof("Removing image %v ...", imageName)

	cmdString := fmt.Sprintf("docker rmi %s", imageName)
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

func (op *ShellOperator) ContainerRun(opts docker.RunOptions, params ...string) error {
	// build the command
	cmdString, err := BuildContainerRunCommand(opts, params...)
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
func BuildContainerRunCommand(opts docker.RunOptions, params ...string) (string, error) {
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
	for _, param := range params {
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

func (op *ShellOperator) ContainerGetPortBinding(container, containerPort, protocol string) (hostPort string, err error) {
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
