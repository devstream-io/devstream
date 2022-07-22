package gitlabcedocker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/pkg/util/docker"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func Create(options map[string]interface{}) (map[string]interface{}, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}

	defaults(&opts)

	if errs := validate(&opts); len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Options error: %s.", e)
		}
		return nil, fmt.Errorf("opts are illegal")
	}

	op := GetDockerOperator(opts)

	// 1. try to pull the image
	// always pull the image because docker will check the image existence
	if err := op.ImagePull(getImageNameWithTag(opts)); err != nil {
		return nil, err
	}

	// 2. try to run the container
	log.Info("Running container as the name <gitlab>")
	if err := op.ContainerRun(buildDockerRunOptions(opts), dockerRunShmSizeParam); err != nil {
		return nil, fmt.Errorf("failed to run container: %v", err)
	}

	// 3. check if the container is started successfully
	if ok := op.ContainerIfRunning(gitlabContainerName); !ok {
		return nil, fmt.Errorf("failed to run container")
	}

	// 4. check if the volume is created successfully
	mounts, err := op.ContainerListMounts(gitlabContainerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get container mounts: %v", err)
	}
	volumes := mounts.ExtractSources()
	if docker.IfVolumesDiffer(volumes, getVolumesDirFromOptions(opts)) {
		return nil, fmt.Errorf("failed to create volumes")
	}

	// 5. show the access url
	showGitLabURL(opts)

	resource := gitlabResource{
		ContainerRunning: true,
		Volumes:          volumes,
		Hostname:         opts.Hostname,
		SSHPort:          strconv.Itoa(int(opts.SSHPort)),
		HTTPPort:         strconv.Itoa(int(opts.HTTPPort)),
		HTTPSPort:        strconv.Itoa(int(opts.HTTPSPort)),
	}

	return resource.toMap(), nil
}

func showGitLabURL(opts Options) {
	accessUrl := opts.Hostname
	if opts.HTTPPort != 80 {
		accessUrl += ":" + strconv.Itoa(int(opts.HTTPPort))
	}
	if !strings.HasPrefix(accessUrl, "http") {
		accessUrl = "http://" + accessUrl
	}

	log.Infof("GitLab access URL: %s", accessUrl)
}
