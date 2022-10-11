package upgrade

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/Masterminds/semver"

	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/interact"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

const (
	assetName      = "dtm-" + runtime.GOOS + "-" + runtime.GOARCH
	dtmTmpFileName = "dtm-tmp"
	dtmBakFileName = "dtm-bak"
	dtmOrg         = "devstream-io"
	dtmRepo        = "devstream"
)

// since dtm file name can be changeable by user,so it should be a variable to get current dtm file name
var dtmFileName string

// Upgrade updates dtm binary file to the latest release version
func Upgrade(continueDirectly bool) error {
	if version.Dev {
		log.Info("Dtm upgrade: do not support to upgrade dtm in development version.")
		os.Exit(0)
	}
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Debugf("Dtm upgrade: work path is : %v", workDir)

	// get dtm bin file path like `/usr/local/bin/dtm-linux-amd64`
	binFilePath, err := os.Executable()
	if err != nil {
		return err
	}
	// get dtm bin file name like `dtm-linux-amd64`
	_, dtmFileName = filepath.Split(binFilePath)

	log.Debugf("Dtm upgrade: dtm file name is : %v", dtmFileName)

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

	// 2. Check whether Upgrade is needed
	// Use Semantic Version to judge. "https://semver.org/"
	shouldUpgrade, err := checkUpgrade(version.Version, ltstReleaseTagName)
	if err != nil {
		return err
	}

	if shouldUpgrade {
		log.Infof("Dtm upgrade: new dtm version: %v is available.", ltstReleaseTagName)
		if !continueDirectly {
			continued := interact.AskUserIfContinue("Would you like to Upgrade? [y/n]")
			if !continued {
				os.Exit(0)
			}
		}

		// 3. Download the latest release version of dtm
		log.Info("Dtm upgrade: downloading the latest release version of dtm, please wait for a while.")
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
// (1) rename current dtm file name to `dtm-bak`.
// (2) rename `dtm-tmp` to current dtm file name.
// (3) grant new dtm file execute permission.
// (4) remove `dtm-bak` binary file.
// TODO(hxcGit): Support for rollback in case of error in intermediate steps
func applyUpgrade(workDir string) error {
	dtmFilePath := filepath.Join(workDir, dtmFileName)
	dtmBakFilePath := filepath.Join(workDir, dtmBakFileName)
	dtmTmpFilePath := filepath.Join(workDir, dtmTmpFileName)

	if err := os.Rename(dtmFilePath, dtmBakFilePath); err != nil {
		return err
	}
	log.Debugf("Dtm upgrade: rename %s to dtm-bak successfully.", dtmFileName)

	if err := os.Rename(dtmTmpFilePath, dtmFilePath); err != nil {
		return err
	}
	log.Debugf("Dtm upgrade: rename dtm-tmp to %s successfully.", dtmFileName)

	if err := os.Chmod(dtmFilePath, 0755); err != nil {
		return err
	}
	log.Debugf("Dtm upgrade: grant %s execute permission successfully.", dtmFileName)

	if err := os.Remove(dtmBakFilePath); err != nil {
		return err
	}
	log.Debug("Dtm upgrade: remove dtm-bak successfully.")

	return nil
}
