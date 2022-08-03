package dockersh

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/docker"
)

var _ = Describe("BuildContainerRunCommand method", func() {
	var opts *docker.RunOptions

	When(" the options are invalid", func() {
		BeforeEach(func() {
			opts = &docker.RunOptions{}
		})

		It("should return an error", func() {
			_, err := BuildContainerRunCommand(opts)
			Expect(err).To(HaveOccurred())
		})
	})

	When(" the options are valid(e.g. gitlab-ce)", func() {
		BeforeEach(func() {
			buildOpts := func() *docker.RunOptions {
				opts := &docker.RunOptions{}
				opts.ImageName = "gitlab/gitlab-ce"
				opts.ImageTag = "rc"
				opts.Hostname = "gitlab.example.com"
				opts.ContainerName = "gitlab"
				opts.RestartAlways = true

				portPublishes := []docker.PortPublish{
					{HostPort: 8122, ContainerPort: 22},
					{HostPort: 8180, ContainerPort: 80},
					{HostPort: 8443, ContainerPort: 443},
				}
				opts.PortPublishes = portPublishes

				gitLabHome := "/srv/gitlab"

				opts.Volumes = []docker.Volume{
					{HostPath: filepath.Join(gitLabHome, "config"), ContainerPath: "/etc/gitlab"},
					{HostPath: filepath.Join(gitLabHome, "data"), ContainerPath: "/var/opt/gitlab"},
					{HostPath: filepath.Join(gitLabHome, "logs"), ContainerPath: "/var/log/gitlab"},
				}

				opts.RunParams = []string{"--shm-size 256m"}

				return opts
			}

			opts = buildOpts()
		})

		It("should return the correct command", func() {
			cmdBuild, err := BuildContainerRunCommand(opts)
			Expect(err).NotTo(HaveOccurred())
			cmdExpect := "docker run --detach --hostname gitlab.example.com" +
				" --publish 8122:22 --publish 8180:80 --publish 8443:443" +
				" --name gitlab --restart always" +
				" --volume /srv/gitlab/config:/etc/gitlab" +
				" --volume /srv/gitlab/data:/var/opt/gitlab" +
				" --volume /srv/gitlab/logs:/var/log/gitlab" +
				" --shm-size 256m gitlab/gitlab-ce:rc"
			Expect(cmdBuild).To(Equal(cmdExpect))
		})

	})
})

var _ = Describe("build[] PortPublish func ", func() {
	var portBindings []string

	BeforeEach(func() {
		portBindings = []string{
			"22/tcp->8122",
			"443/tcp->8443",
			"80/tcp->8180",
		}
	})

	When(" the options are valid", func() {
		It("should return the correct port publishes", func() {
			publishes, err := buildPortPublishes(portBindings)
			Expect(err).NotTo(HaveOccurred())
			Expect(publishes).To(Equal([]docker.PortPublish{
				{HostPort: 8122, ContainerPort: 22},
				{HostPort: 8443, ContainerPort: 443},
				{HostPort: 8180, ContainerPort: 80},
			}))
		})
	})
})
