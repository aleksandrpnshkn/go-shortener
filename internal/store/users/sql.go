package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLStorage struct {
	pgxpool *pgxpool.Pool
}

func (s *SQLStorage) Ping(ctx context.Context) error {
	return s.pgxpool.Ping(ctx)
}

func (s *SQLStorage) Get(ctx context.Context, userID int64) (*User, error) {
	var user User

	row := s.pgxpool.QueryRow(ctx, "SELECT id FROM users WHERE id = $1", userID)
	err := row.Scan(&user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *SQLStorage) Create(ctx context.Context) (*User, error) {
	var userID int64

	row := s.pgxpool.QueryRow(ctx, "INSERT INTO users (id) VALUES (DEFAULT) RETURNING id")
	err := row.Scan(&userID)
	if err != nil {
		return nil, err
	}

	return &User{
		ID: userID,
	}, nil
}

func (s *SQLStorage) Close() error {
	s.pgxpool.Close()
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
