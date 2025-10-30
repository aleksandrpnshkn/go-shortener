package users

import (
	"context"
	"errors"
)

type User struct {
	ID int64
}

type Storage interface {
	Ping(ctx context.Context) error

	Get(ctx context.Context, userID int64) (*User, error)

	Create(ctx context.Context) (*User, error)

	Close() error
}

var (
	ErrUserNotFound = errors.New("user not found")
)
