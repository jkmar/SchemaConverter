package crud_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCrud(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crud Suite")
}
