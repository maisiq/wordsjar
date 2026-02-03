package providers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maisiq/go-words-jar/internal/closer"
	"github.com/maisiq/go-words-jar/internal/config"
	"github.com/maisiq/go-words-jar/internal/container"
	"github.com/maisiq/go-words-jar/internal/logger"
	"github.com/maisiq/go-words-jar/internal/repository/token"
	httpx "github.com/maisiq/go-words-jar/internal/transport/http"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/auth"
	"github.com/maisiq/go-words-jar/internal/transport/http/middleware"
	"github.com/maisiq/go-words-jar/internal/transport/http/router"
)

func IsAdminMiddlewareProvider(c *container.Container) (gin.HandlerFunc, error) {
	authService, err := container.Get[auth.Service](c)
	if err != nil {
		return nil, err
	}

	return middleware.NewIsAdminMiddleware(authService), nil
}

func AuthMiddlewareProvider(c *container.Container) (gin.HandlerFunc, error) {
	secret, err := container.Get[token.SecretRepository](c)
	if err != nil {
		return nil, err
	}

	return middleware.NewAuthMiddleware(secret), nil
}

func RouterProvider(c *container.Container) (*gin.Engine, error) {
	cfg, err := container.Get[config.Config](c)
	if err != nil {
		return nil, err
	}

	// handlers
	auth, err := container.Get[router.AuthHandler](c)
	if err != nil {
		return nil, err
	}

	jar, err := container.Get[router.JarHandler](c)
	if err != nil {
		return nil, err
	}

	words, err := container.Get[router.WordsHandler](c)
	if err != nil {
		return nil, err
	}

	exam, err := container.Get[router.ExamHandler](c)
	if err != nil {
		return nil, err
	}
	handlers := &router.Handlers{
		Auth:  auth,
		Jar:   jar,
		Words: words,
		Exam:  exam,
	}

	// middlewares
	authmw, err := container.GetNamed[gin.HandlerFunc](c, "auth")
	if err != nil {
		return nil, err
	}
	adminmw, err := container.GetNamed[gin.HandlerFunc](c, "admin")
	if err != nil {
		return nil, err
	}

	r := router.New(handlers, &router.Middlewares{
		IsAdmin: adminmw,
		Auth:    authmw,
	}, router.RouterParams{Debug: cfg.App.Debug})
	return r, nil
}

func HTTPServerProvider(c *container.Container) (*httpx.Server, error) {
	cfg, err := container.Get[config.Config](c)
	if err != nil {
		return nil, err
	}

	handler, err := container.Get[*gin.Engine](c)
	if err != nil {
		return nil, err
	}

	log, err := container.Get[logger.Logger](c)
	if err != nil {
		return nil, err
	}

	server := httpx.NewServer(cfg, log, handler)

	closer, err := container.Get[*closer.Closer](c)
	if err != nil {
		return nil, err
	}

	closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.HTTP.ServerShutdownTimeout)*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})

	return server, nil
}
