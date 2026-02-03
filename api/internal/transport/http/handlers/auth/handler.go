package auth

import (
	"context"

	"github.com/maisiq/go-words-jar/internal/config"
	"github.com/maisiq/go-words-jar/internal/models"
	authService "github.com/maisiq/go-words-jar/internal/service/auth"
)

type Service interface {
	CreateUser(ctx context.Context, username, password string) error
	Authenticate(ctx context.Context, username, plainPassword string) (authService.Tokens, error)
	UserInfo(ctx context.Context, username string) (authService.UserInfo, error)
	GetUser(ctx context.Context, username string) (models.User, error)
}

type AuthHandler struct {
	service Service
	cfg     config.HTTP
}

func NewAuthHandler(cfg config.HTTP, service Service) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}
