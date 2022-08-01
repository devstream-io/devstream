package gitlabcedocker

import (
	"strconv"
	"strings"

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

func getGitLabURL(opts *Options) string {
	accessUrl := opts.Hostname
	if opts.HTTPPort != 80 {
		accessUrl += ":" + strconv.Itoa(int(opts.HTTPPort))
	}
	if !strings.HasPrefix(accessUrl, "http") {
		accessUrl = "http://" + accessUrl
	}

	return accessUrl
}

func showGitLabURL(options plugininstaller.RawOptions) error {
	log.Infof("GitLab access URL: %s", gitlabURL)

	return nil
}
