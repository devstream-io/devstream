package devlakeconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func ApplyConfig(options plugininstaller.RawOptions) error {
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
	for _, c := range connections {
		log.Infof("(%s)", c.Name)
		renderAuthConfig(&c)
		configs, err := json.Marshal(c)
		if err != nil {
			return err
		}
		log.Debugf("Connection configs: %s", string(configs))

		url := fmt.Sprintf("%s/plugins/%s/connections", strings.TrimRight(host, "/"), pluginName)
		log.Debugf("URL: %s", url)

		if err := createConnection(url, configs); err != nil {
			return err
		}
	}

	log.Infof("All %s connections have been created.", pluginName)
	return nil
}

func createConnection(url string, bodyWithJson []byte) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyWithJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf(resp.Status)
}

func renderAuthConfig(connection *Connection) {
	connection.Token = connection.Authx.Token
	connection.Username = connection.Authx.Username
	connection.Password = connection.Authx.Password
	connection.AppId = connection.Authx.AppId
	connection.SecretKey = connection.Authx.SecretKey
	connection.Authx = Auth{}
}

func DeleteConfig(options plugininstaller.RawOptions) error {
	// TODO(daniel-hutao): implement later
	return nil
}

func UpdateConfig(options plugininstaller.RawOptions) error {
	// TODO(daniel-hutao): implement later
	return nil
}
