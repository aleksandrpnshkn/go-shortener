package services

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

type Auther interface {
	ParseToken(ctx context.Context, token string) (*users.User, error)

	RegisterUser(ctx context.Context) (*users.User, string, error)

	FromUserContext(ctx context.Context) (*users.User, error)
}

type JwtAuther struct {
	usersStorage users.Storage

	secretKey string
}

type ctxKey string

var (
	ErrInvalidToken = errors.New("invalid token")
)

const ctxUser ctxKey = "user"

func (a *JwtAuther) ParseToken(ctx context.Context, tokenString string) (*users.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims := token.Claims.(*Claims)

	user, err := a.usersStorage.Get(ctx, claims.UserID)
	if err != nil {
		if err == users.ErrUserNotFound {
			return nil, ErrInvalidToken
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (a *JwtAuther) RegisterUser(ctx context.Context) (*users.User, string, error) {
	user, err := a.usersStorage.Create(ctx)
	if err != nil {
		return nil, "", err
	}

	token, err := a.createAuthToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (a *JwtAuther) FromUserContext(ctx context.Context) (*users.User, error) {
	user, ok := ctx.Value(ctxUser).(*users.User)
	if !ok {
		return nil, errors.New("user is not set")
	}

	return user, nil
}

func (a *JwtAuther) createAuthToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           userID,
	})

	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewAuther(usersStorage users.Storage, secretKey string) *JwtAuther {
	return &JwtAuther{
		usersStorage: usersStorage,
		secretKey:    secretKey,
	}
}

func NewUserContext(ctx context.Context, user *users.User) context.Context {
	return context.WithValue(ctx, ctxUser, user)
}
