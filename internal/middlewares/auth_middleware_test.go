package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aleksandrpnshkn/go-shortener/internal/mocks"
	"github.com/aleksandrpnshkn/go-shortener/internal/services"
	"github.com/aleksandrpnshkn/go-shortener/internal/store/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testToken := "testToken"
	testUser := users.User{
		ID: 123,
	}

	t.Run("client sent valid token", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().ParseToken(gomock.Any(), testToken).Return(testUser.ID, nil)
		handler := NewAuthMiddleware(zap.NewExample(), auther, true)(testOkHandler())

		srv := httptest.NewServer(handler)
		defer srv.Close()

		req := httptest.NewRequest("POST", srv.URL, nil)
		req.RequestURI = ""
		req.Header.Add("Cookie", authCookieName+"="+testToken)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "successful status")
	})

	t.Run("client sent invalid token", func(t *testing.T) {
		auther := services.NewAuther(mocks.NewMockUsersStorage(ctrl), "secretkey")
		handler := NewAuthMiddleware(zap.NewExample(), auther, true)(testOkHandler())

		srv := httptest.NewServer(handler)
		defer srv.Close()

		req := httptest.NewRequest("POST", srv.URL, nil)
		req.RequestURI = ""
		req.Header.Add("Cookie", authCookieName+"=blabla")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "unauthorized status")
	})

	t.Run("client not sent token", func(t *testing.T) {
		auther := mocks.NewMockAuther(ctrl)
		auther.EXPECT().RegisterUser(gomock.Any()).Return(testUser.ID, "testToken", nil)
		handler := NewAuthMiddleware(zap.NewExample(), auther, true)(testOkHandler())

		srv := httptest.NewServer(handler)
		defer srv.Close()

		req := httptest.NewRequest("POST", srv.URL, nil)
		req.RequestURI = ""

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode, "successful status")
		assert.Contains(t, resp.Header, "Set-Cookie", "server tries to set cookie")
		assert.Contains(t, resp.Header.Get("Set-Cookie"), testToken, "has token in cookie")
	})
}

func testOkHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	})
}
