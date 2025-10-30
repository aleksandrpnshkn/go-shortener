package services

import (
	"context"
	"errors"
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
	"go.uber.org/zap"
)

type CodesReserver interface {
	GetCode(ctx context.Context) (types.Code, error)
}

type UniqueCodesReserver struct {
	logger        *zap.Logger
	codeGenerator CodeGenerator
	urlsStorage   urls.Storage

	reservedCodes chan types.Code
}

func (u *UniqueCodesReserver) GetCode(ctx context.Context) (types.Code, error) {
	select {
	case <-ctx.Done():
		return types.Code(""), errors.New("context is closed")
	case code, ok := <-u.reservedCodes:
		if !ok {
			return types.Code(""), errors.New("channel for reserved codes is closed")
		}
		return code, nil
	}
}

func (u *UniqueCodesReserver) reserveCodes(ctx context.Context) {
	var code types.Code

	for {
		select {
		case <-ctx.Done():
			u.logger.Info("codes reserver stopped")
			close(u.reservedCodes)
			return
		default:
		}

		code = u.codeGenerator.Generate()

		_, err := u.urlsStorage.Get(ctx, code)
		if err != nil {
			if errors.Is(err, urls.ErrCodeNotFound) {
				u.reservedCodes <- code
				continue
			}

			u.logger.Error("failed to check new code in database, delaying next try...", zap.Error(err))
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (u *UniqueCodesReserver) run(ctx context.Context) {
	go func() {
		u.reserveCodes(ctx)
	}()
}

func NewCodesReserver(
	ctx context.Context,
	logger *zap.Logger,
	codeGenerator CodeGenerator,
	urlsStorage urls.Storage,
	reservedCodesCount int,
) CodesReserver {
	u := &UniqueCodesReserver{
		logger:        logger,
		codeGenerator: codeGenerator,
		urlsStorage:   urlsStorage,

		reservedCodes: make(chan types.Code, reservedCodesCount),
	}

	u.run(ctx)

	return u
}
