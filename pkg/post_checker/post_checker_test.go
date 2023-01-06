package post_checker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Post checker tests", func() {
	It("should not return error for a valid url", func() {
		err := CheckPost("https://github.com")
		Ω(err).Should(BeNil())
	})

	It("should not return error for non-existent url", func() {
		err := CheckPost("https://github.com/no-such-url")
		Ω(err).ShouldNot(BeNil())
	})
})
