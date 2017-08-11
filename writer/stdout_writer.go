package writer

import "fmt"

type StdoutWriter struct {
}

func (stdoutWriter *StdoutWriter) Write(output string) error {
	_, err := fmt.Println(output)
	return err
}
