package trello_test

import (
	. "github.com/onsi/ginkgo/v2"
	//. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

var _ = Describe("Trello", func() {
	Context("Board", func() {
		It("Should does well with board CURD", func() {
			// Adding the your env here
			viper.Set("TRELLO_API_KEY", "")
			viper.Set("TRELLO_TOKEN", "")

			// TODO(daniel-hutao): the code below is only used local now for my TRELLO_API_KEY & TRELLO_TOKEN can't be set in GitHub

			//c, err := trello.NewClient()
			//Expect(err).NotTo(HaveOccurred())
			//err = c.CreateBoard("DS")
			//Expect(err).NotTo(HaveOccurred())
		})
	})
})
