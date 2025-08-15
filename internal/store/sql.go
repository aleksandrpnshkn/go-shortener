package store

import (
	"context"
	"database/sql"

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

func (s *SQLStorage) Set(ctx context.Context, url ShortenedURL) error {
	_, err := s.db.ExecContext(ctx,
		"INSERT INTO urls (code, original_url) VALUES ($1, $2)",
		url.Code,
		url.OriginalURL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLStorage) SetMany(ctx context.Context, urls []ShortenedURL) error {
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO urls (code, original_url) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, url := range urls {
		_, err := stmt.ExecContext(ctx, url.Code, url.OriginalURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLStorage) Get(ctx context.Context, code string) (originalURL string, isFound bool) {
	row := s.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE code = $1", code)
	err := row.Scan(&originalURL)
	if err != nil {
		return "", false
	}

	return originalURL, true
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
