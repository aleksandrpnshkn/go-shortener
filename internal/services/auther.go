package services

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

type Auther interface {
	ParseToken(ctx context.Context, token string) (types.UserID, error)

	RegisterUser(ctx context.Context) (types.UserID, string, error)

	FromUserContext(ctx context.Context) (types.UserID, error)
}

type JwtAuther struct {
	usersStorage users.Storage

	secretKey string
}

type ctxKey string

var (
	ErrInvalidToken = errors.New("invalid token")
)

const ctxUserID ctxKey = "user"

func (a *JwtAuther) ParseToken(ctx context.Context, tokenString string) (types.UserID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.secretKey), nil
	})
	if err != nil {
		return users.GuestID, ErrInvalidToken
	}

	claims := token.Claims.(*Claims)

	return types.UserID(claims.UserID), nil
}

func (a *JwtAuther) RegisterUser(ctx context.Context) (types.UserID, string, error) {
	userID, err := a.usersStorage.Create(ctx)
	if err != nil {
		return users.GuestID, "", err
	}

	token, err := a.createAuthToken(userID)
	if err != nil {
		return users.GuestID, "", err
	}

	return userID, token, nil
}

func (a *JwtAuther) FromUserContext(ctx context.Context) (types.UserID, error) {
	userID, ok := ctx.Value(ctxUserID).(types.UserID)
	if !ok {
		return users.GuestID, errors.New("user is not set")
	}

	return userID, nil
}

func (a *JwtAuther) createAuthToken(userID types.UserID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{},
		UserID:           int64(userID),
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

func NewUserContext(ctx context.Context, userID types.UserID) context.Context {
	return context.WithValue(ctx, ctxUserID, userID)
}
