package urls

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLStorage struct {
	db *sql.DB
}

func (s *SQLStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
}

func (s *SQLStorage) Set(ctx context.Context, url ShortenedURL, user *users.User) (storedURL ShortenedURL, hasConflict bool, err error) {
	const key = "k"
	storedURLs, hasDuplicates, err := s.SetMany(ctx, map[string]ShortenedURL{key: url}, user)
	if err != nil {
		return url, hasDuplicates, err
	}
	return storedURLs[key], hasDuplicates, nil
}

func (s *SQLStorage) SetMany(ctx context.Context, urls map[string]ShortenedURL, user *users.User) (storedURLs map[string]ShortenedURL, hasConflict bool, err error) {
	hasConflict = false

	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO urls (code, original_url, user_id) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (original_url)
		DO UPDATE SET original_url = $2
		RETURNING code, original_url
	`)

	if err != nil {
		return nil, hasConflict, err
	}
	defer stmt.Close()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, hasConflict, err
	}

	storedURLs = make(map[string]ShortenedURL, len(urls))

	for key, url := range urls {
		row := stmt.QueryRowContext(ctx, url.Code, url.OriginalURL, user.ID)

		var storedURL ShortenedURL

		err = row.Scan(&storedURL.Code, &storedURL.OriginalURL)
		if err != nil {
			tx.Rollback()
			return nil, hasConflict, err
		}

		if storedURL.Code != url.Code {
			hasConflict = true
		}

		storedURLs[key] = storedURL
	}

	return storedURLs, hasConflict, tx.Commit()
}

func (s *SQLStorage) Get(ctx context.Context, code string) (originalURL string, err error) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE code = $1", code)
	err = row.Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrCodeNotFound
		}
		return "", err
	}

	return originalURL, nil
}

func (s *SQLStorage) GetByUserID(ctx context.Context, user *users.User) ([]ShortenedURL, error) {
	urls := []ShortenedURL{}

	rows, err := s.db.QueryContext(ctx, "SELECT code, original_url FROM urls WHERE user_id = $1", user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var shortenedURL ShortenedURL

		err = rows.Scan(&shortenedURL.Code, &shortenedURL.OriginalURL)
		if err != nil {
			return nil, err
		}

		urls = append(urls, shortenedURL)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func NewSQLStorage(ctx context.Context, databaseDSN string) (*SQLStorage, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}

	storage := SQLStorage{
		db: db,
	}

	err = storage.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}
