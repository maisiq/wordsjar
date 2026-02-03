package main

import (
	"context"
	"os"

	"github.com/maisiq/go-words-jar/internal/closer"
	"github.com/maisiq/go-words-jar/internal/container"
	"github.com/maisiq/go-words-jar/internal/container/providers"
	httpx "github.com/maisiq/go-words-jar/internal/transport/http"
)

func build() *container.Container {
	c := container.New()

	// common
	container.Provide(c, providers.NewConfigPathProvider("./config/config.yaml")) //TODO: os.Lookup
	container.Provide(c, providers.ConfigProvider)
	container.Provide(c, providers.LoggerProvider)
	container.Provide(c, providers.NewCloserProvider(os.Interrupt, os.Kill))

	// clients
	container.Provide(c, providers.DBClientProvider)

	// storages
	container.Provide(c, providers.RepositoryProvider)
	container.Provide(c, providers.UserRepositoryProvider)
	container.Provide(c, providers.CacheProvider)
	container.Provide(c, providers.SecretRepositoryProvider)

	// services
	container.Provide(c, providers.ServiceProvider)
	container.Provide(c, providers.UserServiceProvider)

	// handlers
	container.Provide(c, providers.AuthHandlerProvider)
	container.Provide(c, providers.ExamHandlerProvider)
	container.Provide(c, providers.JarHandlerProvider)
	container.Provide(c, providers.WordsHandlerProvider)

	// middlewares
	container.ProvideNamed(c, "admin", providers.IsAdminMiddlewareProvider)
	container.ProvideNamed(c, "auth", providers.AuthMiddlewareProvider)

	// http
	container.Provide(c, providers.RouterProvider)
	container.Provide(c, providers.HTTPServerProvider)
	return c
}

func main() {
	ctx := context.Background()
	c := build()

	closer, err := container.Get[*closer.Closer](c)
	if err != nil {
		panic(err)
	}

	server, err := container.Get[*httpx.Server](c)
	if err != nil {
		panic(err)
	}

	server.Run(ctx)
	closer.Wait()
}
