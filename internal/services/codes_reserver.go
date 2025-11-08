package services

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

// CodesReserver резервирует коды для сокращения новых URLов.
type CodesReserver interface {
	GetCode(ctx context.Context) (types.Code, error)
}

// UniqueCodesReserver резервирует уникальные коды.
// Это нужно из-за необходимости проверки на уникальность через БД.
// Для поддержания списка зарезервированных кодов при создании генератора запускается фоновый воркер.
type UniqueCodesReserver struct {
	logger        *zap.Logger
	codeGenerator CodeGenerator
	urlsStorage   urls.Storage

	reservedCodes chan types.Code
}

// GetCode - получить сгенерированный уникальный код.
// Если код ещё не сгенерирован - будет ждать.
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

// NewCodesReserver создаёт генератор кодов
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
