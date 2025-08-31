package searxng_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSearxng(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SearXNG Suite")
}