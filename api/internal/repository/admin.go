package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	errx "github.com/maisiq/go-words-jar/internal/errors"
	"github.com/maisiq/go-words-jar/internal/models"
)

func (repo *Repository) AddWord(ctx context.Context, word models.Word) error {
	query := `INSERT INTO words(id, en, ru, transcription)
	VALUES(:id,:en,:ru,:transcription)
	`
	_, execErr := repo.client.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"id":            word.ID,
		"en":            word.EN,
		"ru":            word.RU,
		"transcription": word.Transcription,
	})

	if execErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(execErr, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return errx.ErrWordAlreadyExists
			}
		}
		return fmt.Errorf("failed to insert user: %w", execErr)
	}

	return nil
}
