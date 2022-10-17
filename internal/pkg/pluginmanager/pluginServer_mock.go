package pluginmanager

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/md5"
)

type mockPluginServer struct {
	*ghttp.Server
}

func newMockPluginServer() *mockPluginServer {
	return &mockPluginServer{Server: ghttp.NewServer()}
}

// registerPluginOK uploads a plugin with the given name and content to the server.
func (s *mockPluginServer) registerPluginOK(plugin, content, wantVersion, os, arch string) {
	version.Version = wantVersion // re-wire version, because many inner functions use version.Version
	tool := configmanager.Tool{Name: plugin}
	pluginFileName := tool.GetPluginFileNameWithOSAndArch(os, arch)
	pluginMD5FileName := tool.GetPluginMD5FileNameWithOSAndArch(os, arch)
	md5Content, _ := md5.CalcMD5(strings.NewReader(content))
	s.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("/v%s/%s", wantVersion, pluginFileName)),
			ghttp.RespondWith(http.StatusOK, content),
		),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("/v%s/%s", wantVersion, pluginMD5FileName)),
			ghttp.RespondWith(http.StatusOK, md5Content),
		),
	)
}

// registerPluginNotFound sets up the server to return a 404 for the given plugin.
func (s *mockPluginServer) registerPluginNotFound(plugin, wantVersion, os, arch string) {
	version.Version = wantVersion // re-wire version, because many inner functions use version.Version
	tool := configmanager.Tool{Name: plugin}
	pluginFileName := tool.GetPluginFileNameWithOSAndArch(os, arch)
	pluginMD5FileName := tool.GetPluginMD5FileNameWithOSAndArch(os, arch)
	s.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("/v%s/%s", wantVersion, pluginFileName)),
			ghttp.RespondWith(http.StatusNotFound, ""),
		),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", fmt.Sprintf("/v%s/%s", wantVersion, pluginMD5FileName)),
			ghttp.RespondWith(http.StatusNotFound, ""),
		),
	)
}
