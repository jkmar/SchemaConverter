package reader

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("reader tests", func() {
	const path = "../tests/"

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

	Describe("read single tests", func() {
		It("Should return error when failed to read a file", func() {
			filename := path + "NonExistingFile.yaml"
			expected := fmt.Errorf(
				"failed to open file %s",
				filename,
			)
			_, err := ReadSingle(filename)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for invalid yaml", func() {
			filename := path + "invalid.yaml"
			expected := fmt.Errorf(
				"cannot parse given schema from file %s",
				filename,
			)
			_, err := ReadSingle(filename)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error when file contains no schemas", func() {
			filename := path + "no_schemas.yaml"
			expected := fmt.Errorf(
				"no schemas found in file %s",
				filename,
			)
			_, err := ReadSingle(filename)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return correct schemas", func() {
			filename := path + "only_names.yaml"
			first := map[interface{}]interface{}{
				"a": "a",
				"b": "b",
			}
			second := map[interface{}]interface{}{
				"c": "c",
				"d": "d",
			}
			expected := []map[interface{}]interface{}{first, second}
			result, err := ReadSingle(filename)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})

	Describe("read all tests", func() {
		It("Should return error when file contains no schemas", func() {
			filename := path + "no_schemas.yaml"
			expected := fmt.Errorf(
				"no schemas found in file %s",
				filename,
			)
			_, err := ReadAll(filename, "")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should return error for file with types other than string", func() {
			filename := path + "no_string_names.yaml"
			expected := fmt.Errorf(
				"in config file %s schemas should be filenames",
				filename,
			)
			_, err := ReadAll(filename, "")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expected))
		})

		It("Should ignore files with no schemas", func() {
			filename := path + "invalid_file_config.yaml"
			expected := []map[interface{}]interface{}{}
			result, err := ReadAll(filename, "")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})

		It("Should return correct schemas", func() {
			filename := path + "only_names_config.yaml"
			first := map[interface{}]interface{}{
				"a": "a",
				"b": "b",
			}
			second := map[interface{}]interface{}{
				"c": "c",
				"d": "d",
			}
			third := map[interface{}]interface{}{
				"e": "e",
				"f": "f",
			}
			fourth := map[interface{}]interface{}{
				"g": "g",
				"h": "h",
			}
			expected := []map[interface{}]interface{}{first, second, third, fourth}
			result, err := ReadAll(filename, "")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})

		It("Should ignore restricted file", func() {
			filename := path + "only_names_config.yaml"
			third := map[interface{}]interface{}{
				"e": "e",
				"f": "f",
			}
			fourth := map[interface{}]interface{}{
				"g": "g",
				"h": "h",
			}
			expected := []map[interface{}]interface{}{third, fourth}
			result, err := ReadAll(filename, path + "only_names.yaml")
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(expected))
		})
	})
})
