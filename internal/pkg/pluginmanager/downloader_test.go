package pluginmanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DownloadClient", Ordered, func() {
	// mock download success func
	mockPlugSuccessGetter := func(reqClient *resty.Client, url, plugName string) error {
		return nil
	}
	// mock download failed func
	mockPlugNotFoundGetter := func(reqClient *resty.Client, url, plugName string) error {
		return fmt.Errorf("downloading plugin %s from %s status code %d", plugName, url, 404)
	}
	var (
		validPlugName    string
		notExistPlugName string
		version          string
		tempDir          string
	)

	BeforeAll(func() {
		tempDir = GinkgoT().TempDir()
		validPlugName = "argocdapp_0.0.1-rc1.so"
		notExistPlugName = "argocdapp_not_exist.so"
		version = "0.0.1-ut-do-not-delete"

	})

	Describe("download method failed", func() {
		var testTable = []struct {
			downloadFunc     func(reqClient *resty.Client, url, plugName string) error
			plugName         string
			expectedErrorMsg string
			describeMsg      string
		}{
			{
				downloadFunc: mockPlugSuccessGetter, plugName: notExistPlugName, expectedErrorMsg: "no such file or directory",
				describeMsg: "should return file not exist if plugin not normal download",
			},
			{
				downloadFunc: mockPlugNotFoundGetter, plugName: validPlugName, expectedErrorMsg: "404",
				describeMsg: "should return 404 if plugin not exist",
			},
		}

		for _, testcase := range testTable {
			It(testcase.describeMsg, func() {
				c := NewDownloadClient()
				c.pluginGetter = testcase.downloadFunc
				err := c.download(tempDir, testcase.plugName, version)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring(testcase.expectedErrorMsg))
			})
		}
	})

	Describe("download method success", func() {
		It("should reanme file if download success", func() {
			tmpFilePath := filepath.Join(tempDir, fmt.Sprintf("%s.tmp", validPlugName))
			f, err := os.Create(tmpFilePath)
			defer os.Remove(tmpFilePath)
			defer f.Close()
			Expect(err).NotTo(HaveOccurred())
			c := NewDownloadClient()
			c.pluginGetter = mockPlugSuccessGetter
			err = c.download(tempDir, validPlugName, version)
			Expect(err).ShouldNot(HaveOccurred())
			renamedFilePath := filepath.Join(tempDir, validPlugName)
			_, err = os.Stat(renamedFilePath)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
