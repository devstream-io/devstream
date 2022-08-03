package gitlabcedocker

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	gitlabImageName       = "gitlab/gitlab-ce"
	defaultImageTag       = "rc"
	gitlabContainerName   = "gitlab"
	dockerRunShmSizeParam = "--shm-size 256m"
)

// gitlabURL is the access URL of GitLab.
var gitlabURL string

func (opts *Options) getGitLabURL() string {
	return fmt.Sprintf("http://%s:%d", opts.Hostname, opts.HTTPPort)
}

func showGitLabURL(options plugininstaller.RawOptions) error {
	log.Infof("GitLab access URL: %s", gitlabURL)

	return nil
}
