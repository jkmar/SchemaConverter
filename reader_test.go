package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("reader tests", func() {
	Describe("get Schemas tests", func() {
		It("Should return error for invalid schema", func() {
			filename := "test"
			object := []interface{}{1, 2}
			expected := fmt.Errorf(
				"error in file %s: schema should have type map[interface{}]interface{}",
				filename,
			)
			_, err := getSchemas(filename, object)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return correct schemas", func() {
			first := map[interface{}]interface{}{
				"a": "a",
				"b": "b",
			}
			second := map[interface{}]interface{}{
				"c": "c",
				"d": "d",
			}
			object := []interface{}{first, second}
			expected := []map[interface{}]interface{}{first, second}
			result, err := getSchemas("", object)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})

	Describe("read schemas from file tests", func() {
		It("Should return error when failed to read a file", func() {
			filename := "tests/NonExistingFile.yaml"
			expected := fmt.Errorf(
				"failed to open file %s",
				filename,
			)
			_, err := ReadSchemasFromFile(filename)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for invalid yaml", func() {
			filename := "tests/invalid.yaml"
			expected := fmt.Errorf(
				"cannot parse given schema from file %s",
				filename,
			)
			_, err := ReadSchemasFromFile(filename)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error when file contains no schemas", func() {
			filename := "tests/no_schemas.yaml"
			expected := fmt.Errorf(
				"no schemas found in file %s",
				filename,
			)
			_, err := ReadSchemasFromFile(filename)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return correct schemas", func() {
			filename := "tests/only_names.yaml"
			first := map[interface{}]interface{}{
				"a": "a",
				"b": "b",
			}
			second := map[interface{}]interface{}{
				"c": "c",
				"d": "d",
			}
			expected := []map[interface{}]interface{}{first, second}
			result, err := ReadSchemasFromFile(filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})
})
