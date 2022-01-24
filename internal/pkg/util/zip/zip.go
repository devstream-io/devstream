package zip

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/merico-dev/stream/internal/pkg/log"
)

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
