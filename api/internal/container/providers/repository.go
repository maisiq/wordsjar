package providers

import (
	"github.com/maisiq/go-words-jar/internal/cache"
	"github.com/maisiq/go-words-jar/internal/container"
	"github.com/maisiq/go-words-jar/internal/db"
	"github.com/maisiq/go-words-jar/internal/repository"
	"github.com/maisiq/go-words-jar/internal/repository/token"
	userRepo "github.com/maisiq/go-words-jar/internal/repository/user"
	"github.com/maisiq/go-words-jar/internal/service"
)

func CacheProvider(c *container.Container) (cache.Cache, error) {
	return cache.NewInMemoryCache(), nil
}

func RepositoryProvider(c *container.Container) (service.Repository, error) {
	db, err := container.Get[*db.DBClient](c)
	if err != nil {
		return nil, err
	}
	repo := repository.New(db)
	return repo, nil
}

func UserRepositoryProvider(c *container.Container) (userRepo.IUserRepository, error) {
	db, err := container.Get[*db.DBClient](c)
	if err != nil {
		return nil, err
	}
	repo := userRepo.NewUserRepository(db)
	return repo, nil
}

func SecretRepositoryProvider(c *container.Container) (token.SecretRepository, error) {
	repo := token.NewInMemory()
	return repo, nil
}
