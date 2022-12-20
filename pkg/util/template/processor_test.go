package template_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/template"
)

var _ = Describe("AddDotForVariablesInConfigProcessor", func() {
	var (
		gotten           []byte
		origin, expected string
	)

	JustBeforeEach(func() {
		gotten = template.AddDotForVariablesInConfigProcessor([]byte(origin))
	})

	When("config is normal", func() {
		BeforeEach(func() {
			origin = "[[varNameA]]"
			expected = "[[ .varNameA ]]"
		})
		It("should succeed", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})

	When("config has spaces", func() {
		BeforeEach(func() {
			origin = "[[ varNameA ]]"
			expected = "[[ .varNameA ]]"
		})

		It("should succeed", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})

	When("config has trailing spaces", func() {
		BeforeEach(func() {
			origin = "[[ varNameA  ]]"
			expected = "[[ .varNameA ]]"
		})

		It("should succeed", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})

	When("config has multiple variables", func() {
		BeforeEach(func() {
			origin = "[[ varNameA ]]/[[ varNameB ]]/[[ varNameC ]]"
			expected = "[[ .varNameA ]]/[[ .varNameB ]]/[[ .varNameC ]]"
		})

		It("should succeed", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})

	When("there are more than one words", func() {
		BeforeEach(func() {
			origin = "[[ func varNameA ]]"
			expected = origin
		})

		It("should do nothing", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})
})

var _ = Describe("AddQuoteForVariablesInConfigProcessor", func() {
	var (
		gotten           []byte
		origin, expected string
	)

	JustBeforeEach(func() {
		gotten = template.AddQuoteForVariablesInConfigProcessor([]byte(origin))
	})

	When("config is normal", func() {
		BeforeEach(func() {
			origin = `[[env GITHUB_TOKEN]]`
			expected = `[[ env "GITHUB_TOKEN" ]]`
		})

		It("should succeed", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})

	When("config has single quote", func() {
		BeforeEach(func() {
			origin = `[[ env 'GITHUB_TOKEN']]`
			expected = origin
		})

		It("should do nothing", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})

	When("config has quote", func() {
		BeforeEach(func() {
			origin = `[[ env "GITHUB_TOKEN"]]`
			expected = origin
		})

		It("should do nothing", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})

	When("config has multiple variables", func() {
		BeforeEach(func() {
			origin = `[[ env GITHUB_TOKEN]]/[[ env "GITLAB_TOKEN"]]`
			expected = `[[ env "GITHUB_TOKEN" ]]/[[ env "GITLAB_TOKEN"]]`
		})

		It("should succeed", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})

	When("the first word has quote", func() {
		BeforeEach(func() {
			origin = `[[ "env" GITHUB_TOKEN]]`
			expected = origin
		})

		It("should do nothing", func() {
			Expect(string(gotten)).To(Equal(expected))
		})
	})

	When("there is only one word", func() {
		BeforeEach(func() {
			origin = `[[GITHUB_TOKEN]]`
			expected = origin
		})

		It("should do nothing", func() {
			log.Infof("------> %s [%s]", string(gotten), expected)
			Expect(string(gotten)).To(Equal(expected))
		})
	})
})
