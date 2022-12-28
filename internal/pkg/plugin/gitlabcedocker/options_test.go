package gitlabcedocker

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	dockerInstaller "github.com/devstream-io/devstream/internal/pkg/plugin/installer/docker"
	"github.com/devstream-io/devstream/pkg/util/docker"
)

var _ = Describe("Options", func() {

	var opts *Options
	BeforeEach(func() {
		opts = &Options{
			GitLabHome:        "/srv/gitlab",
			Hostname:          "gitlab.example.com",
			SSHPort:           8122,
			HTTPPort:          8180,
			HTTPSPort:         8443,
			RmDataAfterDelete: nil,
			ImageTag:          "rc",
		}
	})

	Describe("buildDockerRunOptions func", func() {
		It("should build the docker run options successfully", func() {
			OptsBuild := *buildDockerOptions(opts)
			OptsExpect := dockerInstaller.Options{
				ImageName:     "gitlab/gitlab-ce",
				ImageTag:      "rc",
				Hostname:      "gitlab.example.com",
				ContainerName: "gitlab",
				RestartAlways: true,
				PortPublishes: []docker.PortPublish{
					{HostPort: 8122, ContainerPort: 22},
					{HostPort: 8180, ContainerPort: 80},
					{HostPort: 8443, ContainerPort: 443},
				},
				Volumes: []docker.Volume{
					{HostPath: "/srv/gitlab/config", ContainerPath: "/etc/gitlab"},
					{HostPath: "/srv/gitlab/data", ContainerPath: "/var/opt/gitlab"},
					{HostPath: "/srv/gitlab/logs", ContainerPath: "/var/log/gitlab"},
				},
				RunParams: []string{dockerRunShmSizeParam},
			}

			Expect(OptsBuild).To(Equal(OptsExpect))
		})
	})
})
