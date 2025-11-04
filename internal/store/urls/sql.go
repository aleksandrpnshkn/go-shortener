package urls

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLStorage struct {
	pgxpool *pgxpool.Pool
}

func (s *SQLStorage) Ping(ctx context.Context) error {
	return s.pgxpool.Ping(ctx)
}

func (s *SQLStorage) Close() error {
	s.pgxpool.Close()
	return nil
}

func (s *SQLStorage) Set(ctx context.Context, url ShortenedURL, user *users.User) (storedURL ShortenedURL, hasConflict bool, err error) {
	row := s.pgxpool.QueryRow(ctx, `
		INSERT INTO urls (code, original_url, user_id) 
		VALUES (@code, @originalUrl, @userId) 
		ON CONFLICT (original_url)
		DO UPDATE SET original_url = @originalUrl
		RETURNING code, original_url
	`, pgx.NamedArgs{
		"code":        url.Code,
		"originalUrl": url.OriginalURL,
		"userId":      user.ID,
	})

	err = row.Scan(&storedURL.Code, &storedURL.OriginalURL)
	if err != nil {
		return storedURL, hasConflict, err
	}

	if storedURL.Code != url.Code {
		hasConflict = true
	}

	return storedURL, hasConflict, nil
}

func (s *SQLStorage) SetMany(ctx context.Context, urls map[string]ShortenedURL, user *users.User) (storedURLs map[string]ShortenedURL, hasConflicts bool, err error) {
	hasConflicts = false

	tx, err := s.pgxpool.Begin(ctx)
	if err != nil {
		return nil, hasConflicts, err
	}
	defer tx.Rollback(ctx)

	storedURLs = make(map[string]ShortenedURL, len(urls))

	for key, url := range urls {
		storedURL, hasConflict, err := s.Set(ctx, url, user)
		if err != nil {
			return nil, hasConflicts, err
		}

		if hasConflict {
			hasConflicts = true
		}

		storedURLs[key] = storedURL
	}

	return storedURLs, hasConflicts, tx.Commit(ctx)
}

func (s *SQLStorage) Get(ctx context.Context, code types.Code) (ShortenedURL, error) {
	var shortenedURL ShortenedURL

	row := s.pgxpool.QueryRow(ctx, "SELECT code, original_url, is_deleted FROM urls WHERE code = $1", code)
	err := row.Scan(&shortenedURL.Code, &shortenedURL.OriginalURL, &shortenedURL.IsDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ShortenedURL{}, ErrCodeNotFound
		}
		return ShortenedURL{}, err
	}

	return shortenedURL, nil
}

func (s *SQLStorage) GetByUserID(ctx context.Context, user *users.User) ([]ShortenedURL, error) {
	urls := []ShortenedURL{}

	rows, err := s.pgxpool.Query(ctx, "SELECT code, original_url FROM urls WHERE user_id = $1 AND is_deleted = false", user.ID)
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

func (s *SQLStorage) DeleteManyByUserID(ctx context.Context, commands []DeleteCode) error {
	batch := pgx.Batch{
		QueuedQueries: make([]*pgx.QueuedQuery, 0, len(commands)),
	}

	for _, cmd := range commands {
		batch.Queue("UPDATE urls SET is_deleted = true WHERE user_id = $1 AND code = $2", cmd.User.ID, cmd.Code)
	}

	results := s.pgxpool.SendBatch(ctx, &batch)
	defer results.Close()

	for range commands {
		_, err := results.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewSQLStorage(ctx context.Context, databaseDSN string) (*SQLStorage, error) {
	pool, err := pgxpool.New(ctx, databaseDSN)
	if err != nil {
		return nil, err
	}

	storage := SQLStorage{
		pgxpool: pool,
	}

	err = storage.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}
