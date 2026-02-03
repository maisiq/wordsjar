//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/maisiq/go-words-jar/internal/config"
	dbx "github.com/maisiq/go-words-jar/internal/db"
	"github.com/maisiq/go-words-jar/internal/errors"
	"github.com/maisiq/go-words-jar/internal/repository/user"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_User(t *testing.T) {
	ctx := context.Background()

	username := "testuser"

	client := dbx.New(config.Database{
		DSN: PGDSN,
	})
	repo := user.NewUserRepository(client)

	cases := []struct {
		name     string
		username string
		err      error
	}{
		{"success", username, nil},
		{"non existed username returns error", "non-existed-username", errors.ErrUserNotFound},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			user, err := repo.User(ctx, c.username)

			assert.ErrorIs(t, err, c.err)
			if c.err == nil {
				assert.Equal(t, c.username, user.Username)
			}
		})
	}

}
