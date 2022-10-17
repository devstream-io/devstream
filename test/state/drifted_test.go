package state_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/pluginengine"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/docker"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	dockerUtil "github.com/devstream-io/devstream/pkg/util/docker"
)

var _ = Describe("ResourceDrifted func", func() {
	It("should not be drifted", func() {
		fromFile := resourceFromFile()
		fromRead, err := resourceFromPluginRead()
		Expect(err).To(Succeed())

		drifted, err := pluginengine.ResourceDrifted(fromFile, fromRead)

		Expect(err).To(Succeed())
		Expect(drifted).To(Equal(false))
	})
})

func resourceFromFile() map[string]interface{} {
	configFile := "test_drifted.yaml"
	cfg, err := configmanager.NewManager(configFile).LoadConfig()
	if err != nil {
		panic(err)
	}

	smgr, err := statemanager.NewManager(*cfg.State)
	if err != nil {
		panic(err)
	}

	tool := cfg.Tools[0]

	state := smgr.GetState(statemanager.GenerateStateKeyByToolNameAndInstanceID(tool.Name, tool.InstanceID))

	return state.Resource
}

func resourceFromPluginRead() (map[string]interface{}, error) {
	volumes := []string{
		"/srv/gitlab/config",
		"/srv/gitlab/data",
		"/srv/gitlab/logs",
	}

	portPublishes := []dockerUtil.PortPublish{
		{HostPort: 2022, ContainerPort: 22},
		{HostPort: 8080, ContainerPort: 80},
		{HostPort: 443, ContainerPort: 443},
	}
	resource := &docker.State{
		ContainerRunning: true,
		Volumes:          volumes,
		Hostname:         "gitlab.example.com",
		PortPublishes:    portPublishes,
	}

	return resource.ToMap()
}
