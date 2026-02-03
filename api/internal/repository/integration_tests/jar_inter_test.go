//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/maisiq/go-words-jar/internal/config"
	dbx "github.com/maisiq/go-words-jar/internal/db"
	"github.com/maisiq/go-words-jar/internal/repository"
	"github.com/maisiq/go-words-jar/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetUserWords(t *testing.T) {
	ctx := context.Background()

	username := "testuser"

	client := dbx.New(config.Database{
		DSN: PGDSN,
	})
	repo := repository.New(client)

	noFilters := []service.Filter{}

	ExistedWordID := "73af303e-63c6-4419-8f19-158dbe1f2a3c"
	NotExistedWordID := "73af303e-63c6-4419-8f19-000dbe1f2a3c" // it is possible that testdata could generate this id

	cases := []struct {
		name          string
		username      string
		filters       []service.Filter
		err           error
		wordsQuantity int
	}{
		{"success", username, noFilters, nil, 3},
		{"test mode returns only suitable words", username, []service.Filter{service.WithTestMode()}, nil, 2},
		{"with word id returns only one word", username, []service.Filter{service.WithWordID(ExistedWordID)}, nil, 1},
		{"with wrong word id returns no word", username, []service.Filter{service.WithWordID(NotExistedWordID)}, nil, 0},
		{"with word id && test mode returns no error", username, []service.Filter{service.WithWordID(ExistedWordID)}, nil, 1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			words, err := repo.GetUserWords(ctx, c.username, c.filters...)

			if c.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			assert.Equal(t, c.wordsQuantity, len(words))
		})
	}

}
