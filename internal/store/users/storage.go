package users

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

const GuestID = types.UserID(0)

type User struct {
	ID types.UserID
}

type Storage interface {
	Ping(ctx context.Context) error

	Get(ctx context.Context, userID types.UserID) (*User, error)

	Create(ctx context.Context) (types.UserID, error)

	Close() error
}

var (
	ErrUserNotFound = errors.New("user not found")
)
