package devlakeconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/statemanager"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func ApplyConfig(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}

	for _, p := range opts.Plugins {
		log.Infof("Got DevLake plugin config: %s. Connections: ", p.Name)
		if err := createConnections(opts.DevLakeAddr, p.Name, p.Connections); err != nil {
			return err
		}
	}

	return nil
}

func createConnections(host string, pluginName string, connections []Connection) error {
	for i, c := range connections {
		log.Infof("Connection %d: %s", i, c.Name)
		if err := createConnection(host, pluginName, c); err != nil {
			return err
		}
	}
	log.Infof("All %s connections have been created.", pluginName)
	return nil
}

func createConnection(host string, pluginName string, c Connection) error {
	configs, err := json.Marshal(c)
	if err != nil {
		return err
	}
	log.Debugf("Connection configs: %s", string(configs))
	url := fmt.Sprintf("%s/plugins/%s/connections", strings.TrimRight(host, "/"), pluginName)
	log.Debugf("URL: %s", url)
	if _, err = apiRequest(http.MethodPost, url, configs); err != nil {
		return err
	}
	return nil
}

// update existed connection in devlake backend
func updateConnection(host string, pluginName string, c Connection) error {
	configs, err := json.Marshal(c)
	if err != nil {
		return err
	}
	log.Debugf("UPDATE Connection configs: %s", string(configs))
	url := fmt.Sprintf("%s/plugins/%s/connections/%d", strings.TrimRight(host, "/"), pluginName, c.ID)
	log.Debugf("URL: %s", url)
	if _, err = apiRequest(http.MethodPatch, url, configs); err != nil {
		return err
	}
	return nil
}

// delete existed connection in devlake backend
func deleteConnection(host string, pluginName string, c Connection) error {
	configs, err := json.Marshal(c)
	if err != nil {
		return err
	}
	log.Debugf("DELETE Connection configs: %s", string(configs))
	url := fmt.Sprintf("%s/plugins/%s/connections/%d", strings.TrimRight(host, "/"), pluginName, c.ID)
	log.Debugf("URL: %s", url)
	if _, err = apiRequest(http.MethodDelete, url, configs); err != nil {
		return err
	}
	return nil
}

func getConnections(host string, pluginName string) ([]Connection, error) {
	url := fmt.Sprintf("%s/plugins/%s/connections", strings.TrimRight(host, "/"), pluginName)
	log.Debugf("URL: %s", url)
	resp, err := apiRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	connections := make([]Connection, 0)
	err = json.Unmarshal(resBody, &connections)
	if err != nil {
		return nil, err
	}
	return connections, nil
}

func apiRequest(method string, url string, bodyWithJson []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyWithJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return resp, nil
	}
	return nil, fmt.Errorf(resp.Status)
}

func DeleteConfig(options configmanager.RawOptions) error {
	// TODO(daniel-hutao): implement later
	return nil
}

func UpdateConfig(options configmanager.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	for _, p := range opts.Plugins {
		if updatePluginConfig(opts.DevLakeAddr, p) != nil {
			return err
		}
	}
	return nil
}

func updatePluginConfig(devlakeAddr string, plugin DevLakePlugin) error {
	connections, err := getConnections(devlakeAddr, plugin.Name)
	if err != nil {
		return err
	}

	// Connection.Name -> Connection for config
	connMapForConfig := make(map[string]Connection)
	// Connection.Name -> Connection for status
	connMapForStatus := make(map[string]Connection)
	for _, c := range plugin.Connections {
		connMapForConfig[c.Name] = c
	}
	for _, c := range connections {
		connMapForStatus[c.Name] = c
	}

	return updatePluginConnections(connMapForConfig, connMapForStatus, devlakeAddr, plugin)
}

func updatePluginConnections(connMapForConfig, connMapForStatus map[string]Connection, devlakeAddr string, plugin DevLakePlugin) error {
	for name := range connMapForConfig {
		// Create connection which is not in ResourceStatus
		c, ok := connMapForStatus[name]
		if ok {
			if err := createConnection(devlakeAddr, plugin.Name, connMapForConfig[name]); err != nil {
				return err
			}
			continue
		}

		// Update connection which is different from State
		if c != connMapForConfig[name] {
			if err := updateConnection(devlakeAddr, plugin.Name, connMapForStatus[name]); err != nil {
				return err
			}
		}
	}

	// Delete connection which is not in Config
	for name := range connMapForStatus {
		if _, ok := connMapForConfig[name]; !ok {
			if err := deleteConnection(devlakeAddr, plugin.Name, connMapForStatus[name]); err != nil {
				return err
			}
		}
	}

	return nil
}

func GetStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	// TODO(daniel-hutao): get the real status later
	resStatus := statemanager.ResourceStatus(options)
	return resStatus, nil
}
