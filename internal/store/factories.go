package store

import (
	"context"
	"errors"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
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
