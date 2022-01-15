package githubactions

//Notice: run local only

//import (
//	. "github.com/onsi/ginkgo/v2"
//	. "github.com/onsi/gomega"
//	"gopkg.in/yaml.v3"
//)
//
//var _ = Describe("Githubactions", func() {
//	var optionsStr = `owner: daniel-hutao
//repo: animal
//language:
//  name: go
//  version: "1.17"
//branch: master
//jobs:
//  build:
//    enable: True
//    command: "go build ./..."
//  test:
//    enable: True
//    command: "go test ./..."
//    coverage:
//      enable: False
//      profile: "-race -coverprofile=coverage.out -covermode=atomic"
//`
//	var options *map[string]interface{}
//
//	BeforeEach(func() {
//		  err := yaml.Unmarshal([]byte(optionsStr), &options)
//		  Expect(err).To(Equal(nil))
//	})
//
//	Context("for IsHealthy()", func() {
//		It("should return (true, nil)", func() {
//			healthy, err := IsHealthy(options)
//			Expect(err).To(Equal(nil))
//			Expect(healthy).To(Equal(true))
//		})
//	})
//})
