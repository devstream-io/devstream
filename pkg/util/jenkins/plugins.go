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
	NamePattern = regexp.MustCompile(`^[0-9a-zA-Z\-_]+$`)
)

//go:embed tpl/plugins.tpl.groovy
var pluginsGroovyScript string

func (j *jenkins) validatePluginExist(pluginName string) (bool, error) {
	if ok := NamePattern.MatchString(pluginName); !ok {
		return false, errors.Errorf("invalid plugin name '%s', must follow pattern '%s'", pluginName, NamePattern.String())
	}
	plugin, err := j.HasPlugin(j.ctx, pluginName)
	if err != nil {
		return false, fmt.Errorf("jenkins plugin failed to check plugin %s: %s", pluginName, err)
	}
	if plugin != nil {
		return true, nil
	}
	return false, nil
}

func (j *jenkins) InstallPluginsIfNotExists(plugins []string, enableRestart bool) error {
	toInstallPlugins := make([]string, 0)
	for _, pluginName := range plugins {
		exist, err := j.validatePluginExist(pluginName)
		if err != nil {
			return err
		}
		if !exist {
			log.Debugf("jenkins need to install plugin %s", pluginName)
			toInstallPlugins = append(toInstallPlugins, pluginName)
		} else {
			log.Debugf("jenkins plugin %s already exist", pluginName)
		}
	}
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
	_, err = j.ExecuteScript(pluginInstallScript)
	if err != nil {
		log.Debugf("jenkins install plugins failed:%s", err)
		return err
	}
	// wait jenkins to restart
	if enableRestart {
		tryTime := 6
		return j.waitJenkinsRestart(tryTime)
	} else {
		return errors.New("installed new plugins need to restart jenkins")
	}
}

func (j *jenkins) waitJenkinsRestart(tryTime int) error {
	for {
		status, err := j.Poll(context.TODO())
		if err != nil || status != http.StatusOK {
			log.Debugf("jenkins wait to restart...")
		} else {
			return nil
		}
		time.Sleep(10 * time.Second)
		tryTime = tryTime - 1
		if tryTime == 0 {
			return errors.New("jenkins restart exceed time")
		}
	}
}
