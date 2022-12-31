package post_manager

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPostManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LinkManager Suite")
}
