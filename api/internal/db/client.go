package db

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
	"github.com/maisiq/go-words-jar/internal/config"
	"github.com/maisiq/go-words-jar/internal/logger"
)

type DBClient struct {
	DB *sqlx.DB
}

func New(cfg config.Database) *DBClient {
	db, err := sqlx.Connect("pgx", cfg.DSN)
	if err != nil {
		logger.Fatalw(
			"failed to create db connection",
			"error", err.Error(),
		)
	}

	return &DBClient{
		DB: db,
	}
}
