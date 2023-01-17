package upgrade

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"

	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/interact"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	//the rollback step for the applyUpgrade function
	STEP1 = iota
	STEP2
	STEP3
	STEP4
	COMPLETED

	assetName      = "dtm-" + runtime.GOOS + "-" + runtime.GOARCH
	dtmTmpFileName = "dtm-tmp"
	dtmBakFileName = "dtm-bak"
	dtmDownloadUrl = "https://download.devstream.io/"
	latestVersion  = "https://download.devstream.io/latest_version"
)

// since dtm file name can be changeable by user,so it should be a variable to get current dtm file name
var dtmFileName string

// Upgrade updates dtm binary file to the latest release version
func Upgrade(continueDirectly bool) error {
	if version.Dev {
		log.Info("Dtm upgrade: do not support to upgrade dtm in development version.")
		os.Exit(0)
	}

	// get dtm bin file path like `/usr/local/bin/dtm-linux-amd64`
	binFilePath, err := os.Executable()
	if err != nil {
		return err
	}
	// get dtm bin file name like `dtm-linux-amd64`
	_, dtmFileName = filepath.Split(binFilePath)

	log.Debugf("Dtm upgrade: dtm file name is : %v", dtmFileName)

	workDir := strings.Trim(binFilePath, dtmFileName)
	if err != nil {
		return err
	}
	log.Debugf("Dtm upgrade: work path is : %v", workDir)

	// 1. Get the latest release version
	request, err := http.Get(latestVersion)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(request.Body)
	latestReleaseTagName := string(body[0 : len(body)-1])
	if err != nil {
		return err
	}
	log.Debugf("Dtm upgrade: got latest release version: %v", latestReleaseTagName)

	// 2. Check whether Upgrade is needed
	// To Use Semantic Version to judge. "https://semver.org/"
	shouldUpgrade, err := checkUpgrade(version.Version, latestReleaseTagName)
	if err != nil {
		return err
	}
	if shouldUpgrade {
		log.Infof("Dtm upgrade: new dtm version: %v is available.", latestReleaseTagName)
		if !continueDirectly {
			continued := interact.AskUserIfContinue("Would you like to Upgrade? [y/n]")
			if !continued {
				os.Exit(0)
			}
		}

		// 3. Download the latest release version of dtm
		log.Info("Dtm upgrade: downloading the latest release version of dtm, please wait for a while.")
		var downloadURL strings.Builder
		downloadURL.WriteString(dtmDownloadUrl)
		downloadURL.WriteString(latestReleaseTagName + "/")
		downloadURL.WriteString(assetName)
		downloadSize, downloadError := downloader.New().WithProgressBar().Download(downloadURL.String(), dtmTmpFileName, workDir)
		if downloadError != nil {
			log.Debugf("Failed to download dtm: %v-%v.", latestReleaseTagName, assetName)
			return downloadError
		}
		log.Debugf("Downloaded <%d> bytes.", downloadSize)

		// 4. Replace old dtm with the latest one
		if err = applyUpgrade(workDir); err != nil {
			log.Debug("Failed to replace dtm with latest version.")
			return err
		}
		log.Info("Dtm upgrade successfully!")
		return nil
	}

	log.Info("Dtm upgrade: dtm is the latest version.")
	return nil
}

// checkUpgrade use third-party library `github.com/Masterminds/semver` to compare versions
func checkUpgrade(oldVersion, newVersion string) (bool, error) {
	oldSemVer, newSemVer, err := parseVersion(oldVersion, newVersion)
	if err != nil {
		return false, err
	}

	if oldSemVer.Compare(newSemVer) == -1 {
		return true, nil
	}
	return false, nil
}

// parseVersion parses string version to semver.Version
func parseVersion(old, new string) (*semver.Version, *semver.Version, error) {
	oldSemVer, err := semver.NewVersion(old)
	if err != nil {
		return nil, nil, err
	}

	newSemVer, err := semver.NewVersion(new)
	if err != nil {
		return nil, nil, err
	}

	return oldSemVer, newSemVer, nil
}

// applyUpgrade use os.Rename replace old version dtm with the latest one
// (1) rename current dtm file name to `dtm-bak`.
// (2) rename `dtm-tmp` to current dtm file name.
// (3) grant new dtm file execute permission.
// (4) remove `dtm-bak` binary file.

func applyUpgrade(workDir string) error {
	dtmFilePath := filepath.Join(workDir, dtmFileName)
	dtmBakFilePath := filepath.Join(workDir, dtmBakFileName)
	dtmTmpFilePath := filepath.Join(workDir, dtmTmpFileName)
	updateProgress := STEP1
	defer func() {
		for ; updateProgress >= STEP1; updateProgress-- {
			switch updateProgress {
			//If the error occur when step 1 (rename dtmFileName to `dtm-bak`), delete `dtm-tmp`
			case STEP1:
				if err := os.Remove(dtmTmpFilePath); err != nil {
					log.Debugf("Dtm upgrade rollback error: %s", err.Error())
				}

			//the error occur in the step 2
			case STEP2:
				if err := os.Rename(dtmBakFilePath, dtmFilePath); err != nil {
					log.Debugf("Dtm upgrade rollback error: %s", err.Error())
				}

			//the error occur in the step 3
			case STEP3:
				if err := os.Rename(dtmFilePath, dtmTmpFilePath); err != nil {
					log.Debugf("Dtm upgrade rollback error: %s", err.Error())
				}

			//the error occur in the step 4
			case STEP4:
				if err := os.Chmod(dtmFilePath, 0644); err != nil {
					log.Debugf("Dtm upgrade rollback error: %s", err.Error())
				}
			case COMPLETED:
				//Successfully completed all step
				return
			}
		}
	}()

	if err := os.Rename(dtmFilePath, dtmBakFilePath); err != nil {
		return err
	}
	updateProgress++
	log.Debugf("Dtm upgrade: rename %s to dtm-bak successfully.", dtmFileName)

	if err := os.Rename(dtmTmpFilePath, dtmFilePath); err != nil {
		return err
	}
	updateProgress++
	log.Debugf("Dtm upgrade: rename dtm-tmp to %s successfully.", dtmFileName)

	if err := os.Chmod(dtmFilePath, 0755); err != nil {
		return err
	}
	updateProgress++
	log.Debugf("Dtm upgrade: grant %s execute permission successfully.", dtmFileName)

	if err := os.Remove(dtmBakFilePath); err != nil {
		return err
	}
	updateProgress++
	log.Debug("Dtm upgrade: remove dtm-bak successfully.")

	return nil
}
