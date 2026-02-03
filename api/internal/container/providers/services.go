package providers

import (
	"github.com/maisiq/go-words-jar/internal/cache"
	"github.com/maisiq/go-words-jar/internal/container"
	"github.com/maisiq/go-words-jar/internal/logger"
	"github.com/maisiq/go-words-jar/internal/repository/token"
	userRepo "github.com/maisiq/go-words-jar/internal/repository/user"
	"github.com/maisiq/go-words-jar/internal/service"
	authService "github.com/maisiq/go-words-jar/internal/service/auth"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/auth"
)

func ServiceProvider(c *container.Container) (*service.Service, error) {
	repo, err := container.Get[service.Repository](c)
	if err != nil {
		return nil, err
	}
	log, err := container.Get[logger.Logger](c)
	if err != nil {
		return nil, err
	}
	service := service.New(log, repo)
	return service, nil
}

func UserServiceProvider(c *container.Container) (auth.Service, error) {
	repo, err := container.Get[userRepo.IUserRepository](c)
	if err != nil {
		return nil, err
	}

	storage, err := container.Get[cache.Cache](c)
	if err != nil {
		return nil, err
	}

	log, err := container.Get[logger.Logger](c)
	if err != nil {
		return nil, err
	}

	secret, err := container.Get[token.SecretRepository](c)
	if err != nil {
		return nil, err
	}

	service := authService.NewAuthService(log, repo, storage, secret)
	return service, nil
}
