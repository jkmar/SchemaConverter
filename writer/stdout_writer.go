package writer

import "fmt"

// StdoutWriter is a Writer implementation
type StdoutWriter struct {
}

// Write implementation
func (stdoutWriter *StdoutWriter) Write(output string) error {
	_, err := fmt.Println(output)
	return err
}
