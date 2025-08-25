package users

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLStorage struct {
	db *sql.DB
}

func (s *SQLStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *SQLStorage) Get(ctx context.Context, userID int64) (*User, error) {
	var user User

	row := s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1", userID)
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

	row := s.db.QueryRowContext(ctx, "INSERT INTO users (id) VALUES (DEFAULT) RETURNING id")
	err := row.Scan(&userID)
	if err != nil {
		return nil, err
	}

	return &User{
		ID: userID,
	}, nil
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
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
