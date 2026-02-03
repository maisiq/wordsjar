package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maisiq/go-words-jar/internal/repository/token"
	"github.com/maisiq/go-words-jar/internal/service/auth"
	"github.com/maisiq/go-words-jar/internal/transport/http/middleware/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	ctx := context.Background()

	r := gin.New()
	secretRepo := token.NewInMemory()
	secretMock := mocks.NewSecretServiceMock(t)
	kid := "1"

	t.Run("got valid token from cookie", func(t *testing.T) {
		response := httptest.NewRecorder()
		tctx := gin.CreateTestContextOnly(response, r)
		request := httptest.NewRequest("GET", "/", nil)

		testUsername := "test_user"
		kid, _ = secretRepo.GetKID(ctx)
		unsignedAccessToken, _ := auth.CreateUnsignedAccessToken(auth.AccessTokenInput{
			Username: testUsername,
			KID:      kid,
		})
		token, _ := secretRepo.SignData(ctx, unsignedAccessToken)

		accessCookie := &http.Cookie{
			Name:   "access_token",
			Value:  token,
			MaxAge: int(time.Hour.Seconds()),
		}

		request.AddCookie(accessCookie)
		tctx.Request = request

		mw := NewAuthMiddleware(secretRepo)
		mw(tctx)

		expectedHTTPCode := http.StatusOK

		assert.Equal(t, expectedHTTPCode, response.Code)
	})

	t.Run("got valid token from header", func(t *testing.T) {
		response := httptest.NewRecorder()
		tctx := gin.CreateTestContextOnly(response, r)
		request := httptest.NewRequest("GET", "/", nil)

		testUsername := "test_user"
		kid, _ = secretRepo.GetKID(ctx)
		unsignedAccessToken, _ := auth.CreateUnsignedAccessToken(auth.AccessTokenInput{
			Username: testUsername,
			KID:      kid,
		})
		token, _ := secretRepo.SignData(ctx, unsignedAccessToken)

		request.Header.Set("Authorization", "Bearer "+token)
		tctx.Request = request

		mw := NewAuthMiddleware(secretRepo)
		mw(tctx)

		expectedHTTPCode := http.StatusOK

		assert.Equal(t, expectedHTTPCode, response.Code)
	})

	t.Run("middleware prefers header before cookie", func(t *testing.T) {
		response := httptest.NewRecorder()
		ctx := gin.CreateTestContextOnly(response, r)
		request := httptest.NewRequest("GET", "/", nil)

		testUsername := "test_user"
		kid, _ = secretRepo.GetKID(ctx)
		unsignedAccessToken, _ := auth.CreateUnsignedAccessToken(auth.AccessTokenInput{
			Username: testUsername,
			KID:      kid,
		})
		token, _ := secretRepo.SignData(ctx, unsignedAccessToken)

		// set invalid token in cookie
		accessCookie := &http.Cookie{
			Name:   "access_token",
			Value:  "invalid",
			MaxAge: int(time.Hour.Seconds()),
		}
		request.AddCookie(accessCookie)

		request.Header.Set("Authorization", "Bearer "+token)

		ctx.Request = request

		mw := NewAuthMiddleware(secretRepo)
		mw(ctx)

		expectedHTTPCode := http.StatusOK

		assert.Equal(t, expectedHTTPCode, response.Code)
	})
	t.Run("got invalid token from cookie", func(t *testing.T) {
		response := httptest.NewRecorder()
		tctx := gin.CreateTestContextOnly(response, r)
		request := httptest.NewRequest("GET", "/", nil)
		token := "invalid_token"
		accessCookie := &http.Cookie{
			Name:   "access_token",
			Value:  token,
			MaxAge: int(time.Hour.Seconds()),
		}
		request.AddCookie(accessCookie)
		tctx.Request = request

		mw := NewAuthMiddleware(secretRepo)
		mw(tctx)

		expectedHTTPCode := http.StatusForbidden
		data := response.Body.String()

		assert.Equal(t, expectedHTTPCode, response.Code)
		assert.Equal(t, "{\"detail\":\"invalid token\"}", data)
	})

	t.Run("got invalid token from header", func(t *testing.T) {
		response := httptest.NewRecorder()
		tctx := gin.CreateTestContextOnly(response, r)
		request := httptest.NewRequest("GET", "/", nil)
		token := "invalid_token"
		request.Header.Set("Authorization", "Bearer "+token)
		tctx.Request = request

		mw := NewAuthMiddleware(secretRepo)
		mw(tctx)

		expectedHTTPCode := http.StatusForbidden
		data := response.Body.String()

		assert.Equal(t, expectedHTTPCode, response.Code)
		assert.Equal(t, "{\"detail\":\"invalid token\"}", data)
	})

	t.Run("no token provided", func(t *testing.T) {
		response := httptest.NewRecorder()
		tctx := gin.CreateTestContextOnly(response, r)
		request := httptest.NewRequest("GET", "/", nil)
		tctx.Request = request

		mw := NewAuthMiddleware(secretMock)
		mw(tctx)

		expectedHTTPCode := http.StatusForbidden
		data := response.Body.String()

		assert.Equal(t, expectedHTTPCode, response.Code)
		assert.Equal(t, "{\"detail\":\"no token provided\"}", data)
	})

}
