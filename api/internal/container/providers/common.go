package providers

import (
	"os"

	"github.com/maisiq/go-words-jar/internal/closer"
	"github.com/maisiq/go-words-jar/internal/config"
	"github.com/maisiq/go-words-jar/internal/container"
	"github.com/maisiq/go-words-jar/internal/logger"
)

type ConfigPath string

func NewConfigPathProvider(path string) container.Constructor[ConfigPath] {
	return func(c *container.Container) (ConfigPath, error) {
		return ConfigPath(path), nil
	}
}

func NewCloserProvider(sig ...os.Signal) container.Constructor[*closer.Closer] {
	return func(c *container.Container) (*closer.Closer, error) {
		return closer.New(sig...), nil
	}
}

func ConfigProvider(c *container.Container) (config.Config, error) {
	path, err := container.Get[ConfigPath](c)
	if err != nil {
		return config.Config{}, err
	}
	cfg := config.Load(string(path))
	return cfg, nil
}

func LoggerProvider(c *container.Container) (logger.Logger, error) {
	cfg, err := container.Get[config.Config](c)
	if err != nil {
		return nil, err
	}
	logger.Init(cfg.App.Debug)
	return logger.Get(), nil
}
