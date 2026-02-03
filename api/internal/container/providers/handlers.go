package providers

import (
	"github.com/maisiq/go-words-jar/internal/config"
	"github.com/maisiq/go-words-jar/internal/container"
	"github.com/maisiq/go-words-jar/internal/service"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/auth"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/exam"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/jar"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/words"
	"github.com/maisiq/go-words-jar/internal/transport/http/router"
)

func AuthHandlerProvider(c *container.Container) (router.AuthHandler, error) {
	authService, err := container.Get[auth.Service](c)
	if err != nil {
		return nil, err
	}
	cfg, err := container.Get[config.Config](c)
	if err != nil {
		return nil, err
	}
	handlers := auth.NewAuthHandler(cfg.HTTP, authService)
	return handlers, nil
}

func JarHandlerProvider(c *container.Container) (router.JarHandler, error) {
	service, err := container.Get[*service.Service](c)
	if err != nil {
		return nil, err
	}
	handlers := jar.NewJarHandler(service)
	return handlers, nil
}

func WordsHandlerProvider(c *container.Container) (router.WordsHandler, error) {
	service, err := container.Get[*service.Service](c)
	if err != nil {
		return nil, err
	}
	handlers := words.NewWordsrHandler(service)
	return handlers, nil
}

func ExamHandlerProvider(c *container.Container) (router.ExamHandler, error) {
	service, err := container.Get[*service.Service](c)
	if err != nil {
		return nil, err
	}
	handlers := exam.NewExamHandler(service)
	return handlers, nil
}
