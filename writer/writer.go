package writer

type Writer interface {
	Write(string) error
}

func CreateWriter(filename string) Writer {
	if filename == "" {
		return &StdoutWriter{}
	}
	return &FileWriter{filename}
}
