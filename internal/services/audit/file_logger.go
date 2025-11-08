package audit

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
)

type fileLogger struct {
	file          *os.File
	writer        *bufio.Writer
	lineSeparator rune
}

func (f *fileLogger) addEntry(ctx context.Context, e entry) error {
	return f.addEntries(ctx, []entry{e})
}

func (f *fileLogger) addEntries(ctx context.Context, entries []entry) error {
	lines := [][]byte{}

	for _, entry := range entries {
		line, err := json.Marshal(entry)
		if err != nil {
			return err
		}

		line = append(line, byte(f.lineSeparator))
		lines = append(lines, line)
	}

	err := f.writeLines(lines)
	if err != nil {
		return err
	}

	return nil
}

// Close закрывает файл.
func (f *fileLogger) Close() error {
	err := f.writer.Flush()
	if err != nil {
		return err
	}

	return f.file.Close()
}

func (f *fileLogger) writeLines(lines [][]byte) error {
	for _, line := range lines {
		_, err := f.writer.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewFileLogger создаёт новый файловый логгер для событий аудита
func NewFileLogger(fileName string) (*fileLogger, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}

	fileLogger := &fileLogger{
		file:          file,
		writer:        bufio.NewWriter(file),
		lineSeparator: '\n',
	}

	return fileLogger, nil
}
