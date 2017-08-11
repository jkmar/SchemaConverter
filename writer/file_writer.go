package writer

import "io/ioutil"

type FileWriter struct {
	filename string
}

func (fileWriter *FileWriter) Write(output string) error {
	return ioutil.WriteFile(fileWriter.filename, []byte(output), 0644)
}
