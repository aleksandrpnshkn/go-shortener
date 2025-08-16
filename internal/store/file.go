package store

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
)

type FileStorage struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
	cache   map[string]string

	lastID        int
	lineSeparator rune
}

type ShortenedURLEntry struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (f *FileStorage) Set(ctx context.Context, url ShortenedURL) (storedURL ShortenedURL, hasConflict bool, err error) {
	_, hasConflict, err = f.SetMany(ctx, map[string]ShortenedURL{url.Code: url})
	if err != nil {
		return url, hasConflict, err
	}
	return url, hasConflict, nil
}

func (f *FileStorage) SetMany(ctx context.Context, urls map[string]ShortenedURL) (storedURLs map[string]ShortenedURL, hasConflict bool, err error) {
	lines := [][]byte{}

	for _, url := range urls {
		entry := ShortenedURLEntry{
			UUID:        f.incrementID(),
			ShortURL:    url.Code,
			OriginalURL: url.OriginalURL,
		}

		line, err := json.Marshal(entry)
		if err != nil {
			return nil, false, err
		}

		lines = append(lines, line)
	}

	err = f.writeLines(lines)
	if err != nil {
		return nil, false, err
	}

	for _, url := range urls {
		f.cache[url.Code] = url.OriginalURL
	}

	return urls, false, nil
}

func (f *FileStorage) Get(ctx context.Context, code string) (originalURL string, isFound bool) {
	value, ok := f.cache[code]
	return value, ok
}

func (f *FileStorage) Close() {
	f.file.Close()
}

func (f *FileStorage) incrementID() int {
	f.lastID++
	return f.lastID
}

func (f *FileStorage) writeLines(lines [][]byte) error {
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

		f.cache[entry.ShortURL] = entry.OriginalURL

		if f.lastID < entry.UUID {
			f.lastID = entry.UUID
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

		lastID:        0,
		lineSeparator: '\n',
	}

	err = fileStorage.load()
	if err != nil {
		return nil, err
	}

	return fileStorage, nil
}
