package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
)

type FileStorage struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
	cache   map[string]string

	lastId        int
	lineSeparator rune
}

type ShortenedURLEntry struct {
	UUID        int    `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (f *FileStorage) Set(shortURL string, originalURL string) error {
	entry := ShortenedURLEntry{
		UUID:        f.incrementId(),
		ShortUrl:    shortURL,
		OriginalURL: originalURL,
	}

	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	err = f.writeLine(line)
	if err != nil {
		return err
	}

	f.cache[shortURL] = originalURL
	return nil
}

func (f *FileStorage) Get(shortURL string) (originalURL string, isFound bool) {
	value, ok := f.cache[shortURL]
	return value, ok
}

func (f *FileStorage) Close() {
	f.file.Close()
}

func (f *FileStorage) incrementId() int {
	f.lastId++
	return f.lastId
}

func (f *FileStorage) writeLine(line []byte) error {
	_, err := f.writer.Write(line)
	if err != nil {
		return err
	}

	_, err = f.writer.WriteRune(f.lineSeparator)
	if err != nil {
		return err
	}

	return f.writer.Flush()
}

func (f *FileStorage) readLine() (entry []byte, isFinished bool, err error) {
	if !f.scanner.Scan() {
		err = f.scanner.Err()
		isFinished := err == nil

		return nil, isFinished, err
	}

	entry = f.scanner.Bytes()

	return entry, false, nil
}

func (f *FileStorage) load() error {
	for {
		line, isFinished, err := f.readLine()

		if isFinished {
			return nil
		}

		if err != nil {
			return err
		}

		var entry ShortenedURLEntry
		err = json.Unmarshal(line, &entry)
		if err != nil {
			return errors.New("bad key in file: " + err.Error())
		}

		f.cache[entry.ShortUrl] = entry.OriginalURL

		if f.lastId < entry.UUID {
			f.lastId = entry.UUID
		}
	}
}

func NewFileStorage(fileName string) (*FileStorage, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	fileStorage := &FileStorage{
		file:    file,
		scanner: scanner,
		writer:  bufio.NewWriter(file),
		cache:   map[string]string{},

		lastId:        0,
		lineSeparator: '\n',
	}

	err = fileStorage.load()
	if err != nil {
		return nil, err
	}

	return fileStorage, nil
}
