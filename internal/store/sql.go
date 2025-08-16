package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (s *SQLStorage) Set(ctx context.Context, url ShortenedURL) (storedURL ShortenedURL, hasConflict bool, err error) {
	const key = "k"
	storedURLs, hasDuplicates, err := s.SetMany(ctx, map[string]ShortenedURL{key: url})
	if err != nil {
		return url, hasDuplicates, err
	}
	return storedURLs[key], hasDuplicates, nil
}

func (s *SQLStorage) SetMany(ctx context.Context, urls map[string]ShortenedURL) (storedURLs map[string]ShortenedURL, hasConflict bool, err error) {
	hasConflict = false

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO urls (code, original_url) VALUES ($1, $2)")
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
		_, err := stmt.ExecContext(ctx, url.Code, url.OriginalURL)
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) && pgerr.Code == pgerrcode.UniqueViolation {
			storedCode, isFound := s.getCode(ctx, url.OriginalURL)

			if isFound {
				hasConflict = true
				url.Code = storedCode
				storedURLs[key] = url
				continue
			}
		}
		if err != nil {
			tx.Rollback()
			return nil, hasConflict, err
		}
	}

	return storedURLs, hasConflict, tx.Commit()
}

func (s *SQLStorage) Get(ctx context.Context, code string) (originalURL string, isFound bool) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE code = $1", code)
	err := row.Scan(&originalURL)
	if err != nil {
		return "", false
	}

	return originalURL, true
}

func (s *SQLStorage) getCode(ctx context.Context, originalURL string) (code string, isFound bool) {
	row := s.db.QueryRowContext(ctx, "SELECT code FROM urls WHERE original_url = $1", originalURL)
	err := row.Scan(&code)
	if err != nil {
		return "", false
	}

	return code, true
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
