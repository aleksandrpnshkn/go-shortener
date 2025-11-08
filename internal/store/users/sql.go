package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// SQLStorage - SQL хранилище пользователей.
type SQLStorage struct {
	pgxpool *pgxpool.Pool
}

// Ping проверяет доступность хранилища.
func (s *SQLStorage) Ping(ctx context.Context) error {
	return s.pgxpool.Ping(ctx)
}

// Get достаёт пользователя из хранилища.
func (s *SQLStorage) Get(ctx context.Context, userID types.UserID) (*User, error) {
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

// Create создаёт нового пользователя в хранилище.
func (s *SQLStorage) Create(ctx context.Context) (types.UserID, error) {
	var userID types.UserID

	row := s.pgxpool.QueryRow(ctx, "INSERT INTO users (id) VALUES (DEFAULT) RETURNING id")
	err := row.Scan(&userID)
	if err != nil {
		return GuestID, err
	}

	return userID, nil
}

// Close закрывает соединение с хранилищем.
// Вызывается при завершении работы программы.
func (s *SQLStorage) Close() error {
	s.pgxpool.Close()
	return nil
}

// NewSQLStorage создаёт новое SQL-хранилище пользователей.
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
