package audit

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
)

type FileLogger struct {
	file          *os.File
	writer        *bufio.Writer
	lineSeparator rune
}

func (f *FileLogger) AddEntry(ctx context.Context, entry Entry) error {
	return f.AddEntries(ctx, []Entry{entry})
}

func (f *FileLogger) AddEntries(ctx context.Context, entries []Entry) error {
	lines := [][]byte{}

	for _, entry := range entries {
		line, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		lines = append(lines, line)
	}

	err := f.writeLines(lines)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileLogger) Close() error {
	return f.file.Close()
}

func (f *FileLogger) writeLines(lines [][]byte) error {
	for _, line := range lines {
		_, err := f.writer.Write(line)
		if err != nil {
			return err
		}

		_, err = f.writer.WriteRune(f.lineSeparator)
		if err != nil {
			return err
		}
	}

	return f.writer.Flush()
}

func NewFileLogger(fileName string) (*FileLogger, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}

	fileLogger := &FileLogger{
		file:          file,
		writer:        bufio.NewWriter(file),
		lineSeparator: '\n',
	}

	return fileLogger, nil
}
