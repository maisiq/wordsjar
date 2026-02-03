package providers

import (
	"github.com/maisiq/go-words-jar/internal/closer"
	"github.com/maisiq/go-words-jar/internal/config"
	"github.com/maisiq/go-words-jar/internal/container"
	"github.com/maisiq/go-words-jar/internal/db"
)

func DBClientProvider(c *container.Container) (*db.DBClient, error) {
	cfg, err := container.Get[config.Config](c)
	if err != nil {
		return nil, err
	}

	closer, err := container.Get[*closer.Closer](c)
	if err != nil {
		return nil, err
	}

	client := db.New(cfg.Database)

	closer.Add(func() error {
		return client.DB.Close()
	})

	return client, nil
}
