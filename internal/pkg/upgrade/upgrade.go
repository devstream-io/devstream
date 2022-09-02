package upgrade

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Masterminds/semver"
	"github.com/tcnksm/go-input"

	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

const (
	assetName      = "dtm-" + runtime.GOOS + "-" + runtime.GOARCH
	dtmFileName    = "dtm"
	dtmTmpFileName = "dtm-tmp"
	dtmBakFileName = "dtm-bak"
	dtmOrg         = "devstream-io"
	dtmRepo        = "devstream"
)

// Upgrade updates dtm binary file to the latest release version
func Upgrade(continueDirectly bool) error {
	if version.Dev {
		log.Info("Dtm upgrade: do not support to upgrade dtm in develpment version.")
		os.Exit(0)
	}
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Debugf("Dtm upgrade: work path is : %v", workDir)

	// 1. Get the latest release version
	ghOptions := &git.RepoInfo{
		Org:      dtmOrg,
		Repo:     dtmRepo,
		NeedAuth: false,
		WorkPath: workDir,
	}

	ghClient, err := github.NewClient(ghOptions)
	if err != nil {
		return err
	}

	ltstReleaseTagName, err := ghClient.GetLatestReleaseTagName()
	if err != nil {
		return err
	}
	log.Debugf("Dtm upgrade: got latest release version: %v", ltstReleaseTagName)

	// 2. Chech whether Upgrade is needed
	// Use Semantic Version to judge. "https://semver.org/"
	ok, err := checkUpgrade(version.Version, ltstReleaseTagName)
	if err != nil {
		return err
	}

	if ok {
		log.Infof("Dtm upgrade: new dtm version: %v is available.", ltstReleaseTagName)
		if !continueDirectly {
			userInput := readUserInput()
			if userInput == "n" {
				os.Exit(0)
			}
		}

		// 3. Download the latest release version of dtm
		log.Info("Dtm upgrade: downloading the latest release version of dtm, please wait for a while.")
		// TODO(hxcGit): add download progress bar
		if err = ghClient.DownloadAsset(ltstReleaseTagName, assetName, dtmTmpFileName); err != nil {
			log.Debugf("Failed to download dtm: %v-%v.", ltstReleaseTagName, assetName)
			return err
		}

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
// (1) rename `dtm` to `dtm-bak`.
// (2) rename `dtm-tmp` to `dtm`.
// (3) grant new `dtm` execute permission.
// (4) remove `dtm-bak` binary file.
// TODO(hxcGit): Support for rollback in case of error in intermediate steps
func applyUpgrade(workDir string) error {
	dtmFilePath := filepath.Join(workDir, dtmFileName)
	dtmBakFilePath := filepath.Join(workDir, dtmBakFileName)
	dtmTmpFilePath := filepath.Join(workDir, dtmTmpFileName)

	if err := os.Rename(dtmFilePath, dtmBakFilePath); err != nil {
		return err
	}
	log.Debug("Dtm upgrade: rename dtm to dtm-bak successfully.")

	if err := os.Rename(dtmTmpFilePath, dtmFilePath); err != nil {
		return err
	}
	log.Debug("Dtm upgrade: rename dtm-tmp to dtm successfully.")

	if err := os.Chmod(dtmFilePath, 0755); err != nil {
		return err
	}
	log.Debug("Dtm upgrade: grant dtm execute permission successfully.")

	if err := os.Remove(dtmBakFilePath); err != nil {
		return err
	}
	log.Debug("Dtm upgrade: remove dtm-bak successfully.")

	return nil
}

// TODO(hxcGit): reuse pluginengine.readUserInput()
func readUserInput() string {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "Would you like to Upgrade? [y/n]"
	userInput, err := ui.Ask(query, &input.Options{
		Required: true,
		Default:  "n",
		Loop:     true,
		ValidateFunc: func(s string) error {
			if s != "y" && s != "n" {
				return fmt.Errorf("input must be y or n")
			}
			return nil
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return userInput
}
