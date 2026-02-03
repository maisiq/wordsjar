package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maisiq/go-words-jar/internal/models"
	"github.com/maisiq/go-words-jar/internal/service/auth"
	httpx "github.com/maisiq/go-words-jar/internal/transport/http"
)

const AuthHeaderKey = "Authorization"

func GetUsernameFromContext(ctx *gin.Context) (string, bool) {
	v, ok := ctx.Get(httpx.UsernameContextKey)
	if !ok {
		return "", false
	}
	return v.(string), true
}

type SecretService interface {
	PublicKeys(ctx context.Context) (map[string]any, error)
}

func NewAuthMiddleware(s SecretService) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var err error

		ctx.Header("Vary", AuthHeaderKey)

		// get token from header
		rawAuthHeader := ctx.GetHeader(AuthHeaderKey)
		token, _ := getTokenFromHeader(rawAuthHeader)

		// get token from cookie
		if token == "" {
			token, _ = ctx.Cookie(httpx.AccessTokenKey) // fix
		}

		if token == "" {
			ctx.JSON(http.StatusForbidden, gin.H{"detail": "no token provided"})
			ctx.Abort()
			return
		}

		keys, err := s.PublicKeys(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "iternal error"})
			ctx.Abort()
			return
		}
		data, err := auth.ParseAccessToken(token, keys)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"detail": "invalid token"})
			ctx.Abort()
			return
		}
		ctx.Set(httpx.UsernameContextKey, data.Username)
		ctx.Next()
	}
}

type Service interface {
	GetUser(ctx context.Context, username string) (models.User, error)
}

func NewIsAdminMiddleware(s Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		name, ok := GetUsernameFromContext(ctx)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"detail": "internal error"})
			ctx.Abort()
			return
		}
		user, err := s.GetUser(ctx.Request.Context(), name)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{}) // hide that there is something
			ctx.Abort()
			return
		}
		if !user.IsAdmin {
			ctx.JSON(http.StatusNotFound, gin.H{}) // hide that there is something
			ctx.Abort()
			return
		}

		ctx.Next()

	}

}

func getTokenFromHeader(raw string) (string, error) {
	var err = fmt.Errorf("invalid token")

	s := strings.Split(raw, " ")
	if len(s) != 2 {
		return "", err
	}
	if s[0] != "Bearer" {
		return "", err
	}
	return s[1], nil
}
