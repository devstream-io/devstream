package jenkins

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

var (
	namePattern = regexp.MustCompile(`^[0-9a-zA-Z\-_]+$`)
	// time to wait jenkins restart
	jenkinsRestartRetryTime = 6
)

type JenkinsPlugin struct {
	Name    string
	Version string
}

var basicPlugins = []*JenkinsPlugin{
	{
		Name:    "kubernetes",
		Version: "3600.v144b_cd192ca_a_",
	},
	{
		Name:    "git",
		Version: "4.11.3",
	},
	{
		Name:    "configuration-as-code",
		Version: "1512.vb_79d418d5fc8",
	},
	{
		Name:    "workflow-aggregator",
		Version: "581.v0c46fa_697ffd",
	},
	{
		Name:    "build-user-vars-plugin",
		Version: "1.9",
	},
}

//go:embed tpl/plugins.tpl.groovy
var pluginsGroovyScript string

func (j *jenkins) InstallPluginsIfNotExists(installPlugins []*JenkinsPlugin) error {
	plugins := append(installPlugins, basicPlugins...)
	toInstallPlugins := j.getToInstallPluginList(plugins)
	if len(toInstallPlugins) == 0 {
		return nil
	}

	enableRestart := j.BasicInfo.EnableRestart
	pluginInstallScript, err := template.Render("jenkins-plugins-template", pluginsGroovyScript, map[string]interface{}{
		"JenkinsPlugins": transferPluginSliceToTplString(toInstallPlugins),
		"EnableRestart":  enableRestart,
	})
	if err != nil {
		log.Debugf("jenkins render plugins failed:%s", err)
		return err
	}
	log.Info("jenkins start to install plugins...")
	_, err = j.ExecuteScript(pluginInstallScript)

	// this execute will restart jenkins, so it will return error
	// we just ignore this error to wait for jenkins to restart
	if err != nil {
		log.Debugf("jenkins start restart...")
	}
	log.Debug("jenkins restart to make plugin valid")
	// wait jenkins to restart
	if enableRestart {
		return j.waitJenkinsRestart(toInstallPlugins)
	} else {
		return errors.New("installed new plugins need to restart jenkins")
	}
}

func (j *jenkins) waitJenkinsRestart(toInstallPlugins []*JenkinsPlugin) error {
	tryTime := 1
	for {
		waitTime := tryTime * 20
		// wait 20, 40, 60, 80, 100 seconds for jenkins to restart
		time.Sleep(time.Duration(waitTime) * time.Second)
		log.Infof("wait %d seconds for jenkins plugin install...", waitTime)
		status, err := j.Poll(context.TODO())
		if err == nil && status == http.StatusOK && len(j.getToInstallPluginList(toInstallPlugins)) == 0 {
			return nil
		}
		tryTime++
		if tryTime > jenkinsRestartRetryTime {
			return errors.New("jenkins restart exceed time")
		}
	}
}

func (j *jenkins) getToInstallPluginList(pluginList []*JenkinsPlugin) []*JenkinsPlugin {
	toInstallPlugins := make([]*JenkinsPlugin, 0)
	for _, plugin := range pluginList {
		// check plugin name
		if ok := namePattern.MatchString(plugin.Name); !ok {
			log.Warnf("invalid plugin name '%s', must follow pattern '%s'", plugin.Name, namePattern.String())
			continue
		}
		// check jenkins has this plugin
		installedPlugin, err := j.HasPlugin(j.ctx, plugin.Name)
		if err != nil {
			log.Warnf("jenkins plugin failed to check plugin %s: %s", plugin.Name, err)
			continue
		}

		if installedPlugin == nil {
			log.Debugf("jenkins plugin %s wait to be installed", plugin.Name)
			toInstallPlugins = append(toInstallPlugins, plugin)
		} else {
			log.Debugf("jenkins plugin %s has installed", installedPlugin.ShortName)
		}
	}
	return toInstallPlugins
}

func transferPluginSliceToTplString(plugins []*JenkinsPlugin) string {
	pluginNames := make([]string, 0)
	for _, pluginDetail := range plugins {
		pluginNames = append(pluginNames, fmt.Sprintf("%s:%s", pluginDetail.Name, pluginDetail.Version))
	}
	return strings.Join(pluginNames, ",")
}
