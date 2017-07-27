package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestYamlSchemaToGoStruct(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "YamlSchemaToGoStruct Suite")
}
