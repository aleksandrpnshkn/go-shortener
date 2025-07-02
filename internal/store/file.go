package store

import (
	"bufio"
	"encoding/base64"
	"errors"
	"os"
	"strings"
)

type FileStore struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
	cache   map[string]string

	keyValueSeparator rune
	lineSeparator     rune
}

func (f *FileStore) Set(key string, value string) error {
	// дубли игнорим, на старте загрузится самый актуальный ключ
	line := f.encodeKeyValueLine(key, value)
	err := f.writeLine(line)
	if err != nil {
		return err
	}

	f.cache[key] = value
	return nil
}

func (f *FileStore) Get(key string) (value string, isFound bool) {
	value, ok := f.cache[key]
	return value, ok
}

func (f *FileStore) Close() {
	f.file.Close()
}

func (f *FileStore) writeLine(line string) error {
	_, err := f.writer.WriteString(line)
	if err != nil {
		return err
	}

	err = f.writer.WriteByte(byte(f.lineSeparator))
	if err != nil {
		return err
	}

	return f.writer.Flush()
}

func (f *FileStore) readLine() (keyValue string, isFinished bool, err error) {
	if !f.scanner.Scan() {
		err = f.scanner.Err()
		isFinished := err == nil

		return "", isFinished, err
	}

	keyValue = f.scanner.Text()

	return keyValue, false, nil
}

func (f *FileStore) load() error {
	for {
		line, isFinished, err := f.readLine()

		if isFinished {
			return nil
		}

		if err != nil {
			return err
		}

		key, value, err := f.decodeKeyValueLine(line)
		if err != nil {
			return err
		}

		f.cache[key] = value
	}
}

func (f *FileStore) encodeKeyValueLine(key string, value string) (line string) {
	rawKey := base64.StdEncoding.EncodeToString([]byte(key))
	rawValue := base64.StdEncoding.EncodeToString([]byte(value))

	line = rawKey + string(f.keyValueSeparator) + rawValue

	return line
}

func (f *FileStore) decodeKeyValueLine(line string) (key string, value string, err error) {
	keyValue := strings.Split(line, ":")

	if len(keyValue) != 2 {
		return "", "", errors.New("bad key-value line in file")
	}

	rawKey, err := base64.StdEncoding.DecodeString(keyValue[0])
	if err != nil {
		return "", "", errors.New("bad key in file")
	}
	key = string(rawKey)

	rawValue, err := base64.StdEncoding.DecodeString(keyValue[1])
	if err != nil {
		return "", "", errors.New("bad value in file")
	}
	value = string(rawValue)

	return key, value, nil
}

func NewFileStore(fileName string) (*FileStore, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	fileStore := &FileStore{
		file:    file,
		scanner: scanner,
		writer:  bufio.NewWriter(file),
		cache:   map[string]string{},

		keyValueSeparator: ':',
		lineSeparator:     '\n',
	}

	err = fileStore.load()
	if err != nil {
		return nil, err
	}

	return fileStore, nil
}
