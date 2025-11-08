package urls

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// FileStorage - файловое хранилище сокращённых ссылок.
// Реализует ограниченный набор методов.
type FileStorage struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
	cache   map[types.Code]ShortenedURL

	lastID        int
	lineSeparator rune
}

// ShortenedURLEntry - структура сокращённого URL в файловом хранилище.
type ShortenedURLEntry struct {
	UUID        int               `json:"uuid"`
	Code        types.Code        `json:"short_url"`
	OriginalURL types.OriginalURL `json:"original_url"`
}

// Ping проверяет доступность хранилища.
func (f *FileStorage) Ping(ctx context.Context) error {
	return nil
}

// Set сохраняет короткую ссылку в хранилище.
// Но не проверяет наличие дублей.
func (f *FileStorage) Set(ctx context.Context, url ShortenedURL, userID types.UserID) (storedURL ShortenedURL, hasConflict bool, err error) {
	_, hasConflict, err = f.SetMany(ctx, map[string]ShortenedURL{string(url.Code): url}, userID)
	if err != nil {
		return url, hasConflict, err
	}
	return url, hasConflict, nil
}

// SetMany сохраняет множество коротких ссылок в хранилищк.
// Но не проверяет проверяет наличие дублей.
func (f *FileStorage) SetMany(ctx context.Context, urls map[string]ShortenedURL, userID types.UserID) (storedURLs map[string]ShortenedURL, hasConflicts bool, err error) {
	lines := [][]byte{}

	for _, url := range urls {
		entry := ShortenedURLEntry{
			UUID:        f.incrementID(),
			Code:        url.Code,
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
		f.cache[url.Code] = url
	}

	return urls, false, nil
}

// Get достаёт сокращённую ссылку по её коду.
func (f *FileStorage) Get(ctx context.Context, code types.Code) (ShortenedURL, error) {
	value, ok := f.cache[code]
	if !ok {
		return ShortenedURL{}, ErrCodeNotFound
	}
	return value, nil
}

// GetByUserID - заглушка для доставания всех сокращённых ссылок пользователя.
func (f *FileStorage) GetByUserID(ctx context.Context, userID types.UserID) ([]ShortenedURL, error) {
	return []ShortenedURL{}, nil
}

// DeleteManyByUserID - заглушка для удаления всех сокращённых ссылок пользователя.
func (f *FileStorage) DeleteManyByUserID(ctx context.Context, commands []DeleteCode) error {
	return nil
}

// Close закрывает файл.
func (f *FileStorage) Close() error {
	return f.file.Close()
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

		f.cache[entry.Code] = ShortenedURL{
			OriginalURL: entry.OriginalURL,
			Code:        entry.Code,
			IsDeleted:   false,
		}

		if f.lastID < entry.UUID {
			f.lastID = entry.UUID
		}
	}
}

// NewFileStorage создаёт новое файловое хранилище для URLов.
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
		cache:   map[types.Code]ShortenedURL{},

		lastID:        0,
		lineSeparator: '\n',
	}

	err = fileStorage.load()
	if err != nil {
		return nil, err
	}

	return fileStorage, nil
}
