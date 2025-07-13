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
