package writer

// Writer is an interface for displaying output
type Writer interface {
	Write(string) error
}

// CreateWriter creates a Writer that writes to a file with given filename
func CreateWriter(filename string) Writer {
	if filename == "" {
		return &StdoutWriter{}
	}
	return &FileWriter{filename}
}
