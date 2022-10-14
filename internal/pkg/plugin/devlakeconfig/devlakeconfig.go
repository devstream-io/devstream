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
	// TODO(daniel-hutao): implement later
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	for _, p := range opts.Plugins {
		connections, err := getConnections(opts.DevLakeAddr, p.Name)
		if err != nil {
			return err
		}
		// Map Connection.Name -> Connection for config
		configConnectionMap := make(map[string]Connection)
		// Map Connection.Name -> Connection for state
		stateConnectionMap := make(map[string]Connection)
		for _, configConnection := range p.Connections {
			configConnectionMap[configConnection.Name] = configConnection
		}
		// Construct a map
		for _, stateConnection := range connections {
			stateConnectionMap[stateConnection.Name] = stateConnection
		}
		for k := range configConnectionMap {
			// Create connection which is not in State
			if _, ok := stateConnectionMap[k]; !ok {
				if err = createConnection(opts.DevLakeAddr, p.Name, configConnectionMap[k]); err != nil {
					return err
				}
				continue
			}
			// Update connection which is different from State
			if stateConnectionMap[k] != configConnectionMap[k] {
				if err = updateConnection(opts.DevLakeAddr, p.Name, stateConnectionMap[k]); err != nil {
					return err
				}
				continue
			}
		}
		for k := range stateConnectionMap {
			// Delete connection which is not in config
			if _, ok := configConnectionMap[k]; !ok {
				if err = deleteConnection(opts.DevLakeAddr, p.Name, stateConnectionMap[k]); err != nil {
					return err
				}
				continue
			}
		}
	}
	return nil
}

func GetStatus(options configmanager.RawOptions) (statemanager.ResourceStatus, error) {
	resStatus := statemanager.ResourceStatus(options)
	return resStatus, nil
}
