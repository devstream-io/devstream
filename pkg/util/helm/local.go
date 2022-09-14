package helm

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/md5"
)

func cacheChartPackage(chartPath string) error {
	if chartPath == "" {
		log.Debugf("The option chartPath == \"\".")
		return nil
	}

	// source file checks
	sFile, err := os.Stat(chartPath)
	if err != nil {
		return fmt.Errorf("chart <%s> doesn't exist", chartPath)
	}
	if !sFile.Mode().IsRegular() {
		return fmt.Errorf("the chart file <%s> is non-regular (%q)", chartPath, sFile.Mode().String())
	}

	// destination chart path checks
	dFilePath := filepath.Join(repositoryCache, chartPath)
	log.Debugf("The destination chart path is <%s>.", dFilePath)
	dFile, err := os.Stat(dFilePath)
	if err != nil { // return err if err != nil and err != "NotExist"
		if !os.IsNotExist(err) {
			log.Errorf("Got error: %s.", err)
			return err
		}
	} else { // err == nil -> file exists -> check if the source file equals destination file
		if !(dFile.Mode().IsRegular()) {
			err = fmt.Errorf("the destination file <%s> is non-regular (%q)", dFilePath, dFile.Mode().String())
			log.Errorf("Got error: %s.", err)
			return err
		}
		if os.SameFile(sFile, dFile) {
			log.Debugf("The source chart package and destination chart package are same.")
			return nil
		}
		// check md5
		equal, err := md5.FilesMD5Equal(chartPath, dFilePath)
		if err != nil {
			log.Errorf("Got error: %s.", err)
			return err
		}
		if equal {
			log.Infof("The chart package already exists in the cache directory.")
			return nil
		}
		// remove the destination file if its name equals to source file but their contents don't equal
		if err = os.RemoveAll(dFilePath); err != nil {
			log.Errorf("Got error: %s.", err)
			return err
		}
	}

	// destination chart path is empty, then create it.
	log.Debugf("Prepare to copy <%s> to <%s>.", chartPath, dFilePath)
	if err = os.MkdirAll(repositoryCache, 0755); err != nil {
		return err
	}
	return copyFile(chartPath, dFilePath)
}

func copyFile(srcFile, dstFile string) (err error) {
	// prepare source file
	sFile, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := sFile.Close(); closeErr != nil {
			log.Errorf("Failed to close file %s: %s", srcFile, closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// create destination file
	dFile, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dFile.Close(); closeErr != nil {
			log.Errorf("Failed to close file %s: %s", dstFile, closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// copy and sync
	if _, err = io.Copy(dFile, sFile); err != nil {
		return nil
	}
	return dFile.Sync()
}
