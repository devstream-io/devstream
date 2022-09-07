package jenkins

import (
	"context"
	_ "embed"
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

//go:embed tpl/plugins.tpl.groovy
var pluginsGroovyScript string

func (j *jenkins) InstallPluginsIfNotExists(plugins []string, enableRestart bool) error {
	toInstallPlugins := j.getToInstallPluginList(plugins)
	if len(toInstallPlugins) == 0 {
		return nil
	}

	pluginInstallScript, err := template.Render("jenkins-plugins-template", pluginsGroovyScript, map[string]interface{}{
		"JenkinsPlugins": strings.Join(toInstallPlugins, ","),
		"EnableRestart":  enableRestart,
	})
	if err != nil {
		log.Debugf("jenkins render plugins failed:%s", err)
		return err
	}
	log.Debug("jenkins start to install plugins...")
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

func (j *jenkins) waitJenkinsRestart(toInstallPlugins []string) error {
	tryTime := 1
	for {
		waitTime := tryTime * 20
		// wait 20, 40, 60, 80, 100 seconds for jenkins to restart
		time.Sleep(time.Duration(waitTime) * time.Second)
		log.Debugf("wait %d seconds for jenkins plugin install...", waitTime)
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

func (j *jenkins) getToInstallPluginList(pluginList []string) []string {
	toInstallPlugins := make([]string, 0)
	for _, pluginName := range pluginList {
		// check plugin name
		if ok := namePattern.MatchString(pluginName); !ok {
			log.Warnf("invalid plugin name '%s', must follow pattern '%s'", pluginName, namePattern.String())
			continue
		}
		// check jenkins has this plugin
		plugin, err := j.HasPlugin(j.ctx, pluginName)
		if err != nil {
			log.Warnf("jenkins plugin failed to check plugin %s: %s", pluginName, err)
			continue
		}

		if plugin == nil {
			log.Debugf("jenkins plugin %s wait to be installed", pluginName)
			toInstallPlugins = append(toInstallPlugins, pluginName)
		} else {
			log.Debugf("jenkins plugin %s has installed", pluginName)
		}
	}
	return toInstallPlugins
}
