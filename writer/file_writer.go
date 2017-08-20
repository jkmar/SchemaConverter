package writer

import "io/ioutil"

// FileWriter is a writer implementation
type FileWriter struct {
	filename string
}

// Write implementation
func (fileWriter *FileWriter) Write(output string) error {
	return ioutil.WriteFile(fileWriter.filename, []byte(output), 0644)
}
