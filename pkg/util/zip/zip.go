package zip

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func UnZip(filePath, targetPath string) error {
	archive, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}

	if err = handleArchiveFiles(targetPath, archive.File); err != nil {
		return err
	}

	return archive.Close()
}

func handleArchiveFiles(targetPath string, files []*zip.File) error {
	for _, f := range files {
		filePath := filepath.Join(targetPath, f.Name)
		log.Debugf("Unzipping file -> %s.", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(targetPath)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path")
		}

		if f.FileInfo().IsDir() {
			log.Debugf("Creating directory -> %s.", f.FileInfo().Name())
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err = io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		if err = dstFile.Close(); err != nil {
			return err
		}

		if err = fileInArchive.Close(); err != nil {
			return err
		}
	}
	return nil
}

func UnTargz(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if strings.Contains(header.Name, "._") {
			continue
		}
		if err := handleHeader(header, tarReader); err != nil {
			return err
		}
	}
	return nil
}

func handleHeader(header *tar.Header, reader *tar.Reader) error {
	switch header.Typeflag {
	case tar.TypeDir:
		if err := os.MkdirAll(header.Name, 0755); err != nil {
			return err
		}
	case tar.TypeReg:
		outFile, err := os.Create(header.Name)
		if err != nil {
			return err
		}
		if _, err := io.Copy(outFile, reader); err != nil {
			return err
		}
		if err := outFile.Close(); err != nil {
			return err
		}

	default:
		errMsg := fmt.Sprintf("got unknown type: %b in %s", header.Typeflag, header.Name)
		log.Error(errMsg)
		return fmt.Errorf(errMsg)
	}
	return nil
}
