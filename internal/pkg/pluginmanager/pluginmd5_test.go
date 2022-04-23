package pluginmanager

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/pkg/util/md5"
)

var _ = Describe("CheckLocalPlugins", func() {
	var err error
	var config *configloader.Config
	var file, fileMD5 string
	var tools []configloader.Tool

	Context("MD5", func() {
		BeforeEach(func() {
			viper.Set("plugin-dir", "./")
			tools = []configloader.Tool{
				{InstanceID: "a", Name: "a"},
			}
			config = &configloader.Config{Tools: tools}

			file = configloader.GetPluginFileName(&tools[0])
			fileMD5 = configloader.GetPluginMD5FileName(&tools[0])
			err := createNewFile(file)
			Expect(err).NotTo(HaveOccurred())
		})

		It("CheckLocalPlugins, md5 sum should match .md5 file content", func() {
			err = addMD5File(file, fileMD5)
			Expect(err).NotTo(HaveOccurred())

			err = CheckLocalPlugins(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("CheckLocalPlugins, md5 sum should mismatch .md5 file content", func() {
			err = createNewFile(fileMD5)
			Expect(err).NotTo(HaveOccurred())

			err = CheckLocalPlugins(config)
			expectErrMsg := fmt.Sprintf("plugin %s doesn't match with .md5", tools[0].InstanceID)
			Expect(err.Error()).To(Equal(expectErrMsg))
		})

		It("pluginAndMD5Matches, md5 sum should match .md5 file content", func() {
			err = addMD5File(file, fileMD5)
			Expect(err).NotTo(HaveOccurred())

			err = pluginAndMD5Matches(viper.GetString("plugin-dir"), file, fileMD5, tools[0].InstanceID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("pluginAndMD5Matches, md5 sum should mismatch .md5 file content", func() {
			err = createNewFile(fileMD5)
			Expect(err).NotTo(HaveOccurred())

			err = pluginAndMD5Matches(viper.GetString("plugin-dir"), file, fileMD5, tools[0].InstanceID)
			expectErrMsg := fmt.Sprintf("plugin %s doesn't match with .md5", tools[0].InstanceID)
			Expect(err.Error()).To(Equal(expectErrMsg))
		})

		AfterEach(func() {
			err = os.Remove(file)
			Expect(err).NotTo(HaveOccurred())
			err = os.Remove(fileMD5)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

func createNewFile(fileName string) error {
	f1, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f1.Close()
	return nil
}

func addMD5File(fileName, md5FileName string) error {
	md5, err := md5.CalcFileMD5(fileName)
	if err != nil {
		return err
	}
	md5File, err := os.Create(md5FileName)
	if err != nil {
		return err
	}
	_, err = md5File.Write([]byte(md5))
	if err != nil {
		return err
	}
	return nil
}
