package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("parseNestedVars", func() {
	var (
		origin, expected, parsed map[string]any
		err                      error
	)

	JustBeforeEach(func() {
		parsed, err = parseNestedVars(origin)
	})

	When("vars map is correct", func() {
		When("case is simple(no nested)", func() {
			BeforeEach(func() {
				origin = map[string]any{
					"a": "a",
					"b": 123,
				}
				expected = origin
			})
			It("should parse succeed", func() {
				Expect(err).Should(Succeed())
				Expect(parsed).Should(Equal(expected))
			})
		})
		When("case is complex(nested once)", func() {
			BeforeEach(func() {
				origin = map[string]any{
					"a": "[[ b]]a",
					"b": "123",
				}
				expected = map[string]any{
					"a": "123a",
					"b": "123",
				}
			})
			It("should parse succeed", func() {
				Expect(err).Should(Succeed())
				Expect(parsed).Should(Equal(expected))
			})
		})
		When("case is complex(nested many times)", func() {
			BeforeEach(func() {
				origin = map[string]any{
					"a": "[[ b]]a",
					"b": 123,
					"c": "[[a]]c",
					"d": "[[a]]/[[ c ]]/[[b]]",
				}
				expected = map[string]any{
					"a": "123a",
					"b": 123,
					"c": "123ac",
					"d": "123a/123ac/123",
				}
			})

			It("should parse succeed", func() {
				Expect(err).Should(Succeed())
				Expect(parsed).Should(Equal(expected))
			})
		})
	})

	When("vars map is incorrect", func() {
		BeforeEach(func() {
			origin = map[string]any{
				"a": "[[ b]]a",
				"b": "123",
				"c": "[[a]]c[[c]]",
				"d": "[[a]]/[[ c ]]/[[b]]",
			}
			expected = map[string]any{
				"a": "123a",
				"b": "123",
				"c": "123ac",
				"d": "123a/123ac/123",
			}
		})

		It("should return error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})
})

var _ = Describe("ifContainVar", func() {
	var (
		value       string
		checkResult bool
	)

	JustBeforeEach(func() {
		checkResult = ifContainVar(value)
	})

	When("contains var", func() {
		AfterEach(func() {
			Expect(checkResult).To(Equal(true))
		})
		Context("case 1", func() {
			BeforeEach(func() {
				value = " [[a]]"
			})
			It("should check correctly", func() {})
		})
		Context("case 2", func() {
			BeforeEach(func() {
				value = "[[ b ]]"
			})
			It("should check correctly", func() {})
		})
	})

	When("doesn't contain var", func() {
		AfterEach(func() {
			Expect(checkResult).To(Equal(false))
		})
		Context("case 1", func() {
			BeforeEach(func() {
				value = " [[a] ]"
			})
			It("should check correctly", func() {})
		})
		Context("case 2", func() {
			BeforeEach(func() {
				value = " [ b ]]"
			})
			It("should check correctly", func() {})
		})
	})
})
