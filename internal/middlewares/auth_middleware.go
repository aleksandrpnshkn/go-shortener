package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/types"
)

const authCookieName = "auth_token"

func NewAuthMiddleware(
	logger *zap.Logger,
	auther services.Auther,
	enableRegistration bool,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			authCookie, err := req.Cookie(authCookieName)
			if err != nil && err != http.ErrNoCookie {
				logger.Error("unknown cookie error", zap.Error(err))
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			var token string
			var userID types.UserID

			if err == http.ErrNoCookie {
				if enableRegistration {
					userID, token, err = auther.RegisterUser(ctx)
					if err != nil {
						logger.Error("failed to register new user", zap.Error(err))
						res.WriteHeader(http.StatusInternalServerError)
						return
					}

					authCookie = &http.Cookie{
						Name:  authCookieName,
						Value: token,

						HttpOnly: true,
						SameSite: http.SameSiteStrictMode,
						Secure:   false,
					}
					http.SetCookie(res, authCookie)
				}
			} else {
				userID, err = auther.ParseToken(ctx, authCookie.Value)
				if err != nil {
					if err == services.ErrInvalidToken {
						res.WriteHeader(http.StatusUnauthorized)
						return
					} else {
						logger.Error("failed to register new user", zap.Error(err))
						res.WriteHeader(http.StatusInternalServerError)
						return
					}
				}
			}

			req = req.WithContext(services.NewUserContext(ctx, userID))

			next.ServeHTTP(res, req)
		})
	}
}
