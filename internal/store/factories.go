package store

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"go.uber.org/zap"
)

func NewURLsStorage(
	ctx context.Context,
	databaseDSN string,
	fileStoragePath string,
	logger *zap.Logger,
) (urls.Storage, error) {
	var storage urls.Storage

	storage, err := urls.NewSQLStorage(ctx, databaseDSN)
	if err == nil {
		err = runMigrations(databaseDSN)
		if err != nil {
			return nil, errors.New("failed to run SQL migrations")
		}

		return storage, nil
	}

	logger.Warn("failed to init urls SQL storage", zap.Error(err))

	storage, err = urls.NewFileStorage(fileStoragePath)
	if err != nil {
		return nil, errors.New("failed to init urls file storage")
	}

	return storage, nil
}

func NewUsersStorage(
	ctx context.Context,
	databaseDSN string,
	logger *zap.Logger,
) (users.Storage, error) {
	var storage users.Storage

	storage, err := users.NewSQLStorage(ctx, databaseDSN)
	if err != nil {
		logger.Warn("failed to init users SQL storage", zap.Error(err))
		storage = users.NewMemoryStorage()
	}

	return storage, nil
}
