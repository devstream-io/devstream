package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/merico-dev/stream/cmd/devstream/version"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/pluginmanager"
	"github.com/merico-dev/stream/pkg/util/log"
)

const DTMRemoteMD5 = ".remote"

var initCMD = &cobra.Command{
	Use:   "init",
	Short: "Download needed plugins according to the config file",
	Long:  `Download needed plugins according to the config file`,
	Run:   initCMDFunc,
}

func initCMDFunc(cmd *cobra.Command, args []string) {
	cfg := configloader.LoadConf(configFile)
	if cfg == nil {
		log.Fatal("Failed to load the config file.")
	}

	log.Info("Initialize started.")
	if err := compareDtmMD5(); err != nil {
		log.Error(err)
		return
	}

	err := pluginmanager.DownloadPlugins(cfg)
	if err != nil {
		log.Errorf("Error: %s.", err)
		return
	}

	log.Success("Initialize finished.")
}

// compareDtmMD5 compare dtm itself and remote md5, if does not match, cannot  download plugins from remote
func compareDtmMD5() error {
	// get remote dtm .md5
	dtmMD5FileName := configloader.GetDtmMD5FileName()
	if err := pluginmanager.DownloadDtmMD5(DTMRemoteMD5, dtmMD5FileName); err != nil {
		return err
	}

	selfDtmName, err := getDtmSelf()
	if err != nil {
		return err
	}

	isMD5Match, err := version.ValidateFileMatchMD5(selfDtmName, filepath.Join(DTMRemoteMD5, dtmMD5FileName))
	if err != nil {
		return err
	}
	if !isMD5Match {
		return fmt.Errorf("DTM %s doesn't match with remote .md5, cannot download Plugins from remote release version", selfDtmName)
	}
	return nil
}

// getDtmSelf get executing dtm itself
func getDtmSelf() (string, error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return path, nil
}
