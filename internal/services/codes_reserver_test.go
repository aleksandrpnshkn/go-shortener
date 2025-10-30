package services

import (
	"context"
	"testing"
	"time"

	"github.com/aleksandrpnshkn/go-shortener/internal/mocks"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/urls"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestCodesReserver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("generates unique url", func(t *testing.T) {
		codeLength := 8
		codeGenerator := NewRandomCodeGenerator(codeLength)

		urlsStorage := mocks.NewMockURLsStorage(ctrl)
		urlsStorage.EXPECT().Get(gomock.Any(), gomock.Any()).
			AnyTimes().
			Return(urls.ShortenedURL{}, urls.ErrCodeNotFound)

		codesReserver := NewCodesReserver(
			context.Background(),
			zap.NewExample(),
			codeGenerator,
			urlsStorage,
			1,
		)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		code, err := codesReserver.GetCode(ctx)

		assert.NoError(t, err)
		assert.Len(t, code, codeLength, "should return valid code")
	})
}
