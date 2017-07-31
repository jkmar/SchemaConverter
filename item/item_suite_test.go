package item_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestItem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Item Suite")
}
