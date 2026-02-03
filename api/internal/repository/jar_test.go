package repository

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	dbx "github.com/maisiq/go-words-jar/internal/db"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestRepository_AddWordToJar(t *testing.T) {
	ctx := context.Background()
	username := "test_user"
	words := []string{"word1", "word2"}
	rating := float32(0)
	attempts := 1
	lastAttempt := AnyTime{}

	db, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual),
	)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expectQuery := `WITH subq as NOT MATERIALIZED (
		SELECT en FROM words WHERE en IN ($1, $2) 
	)
	INSERT INTO user_words (username, word_en, knowledge_rating, consecutive_success_attempts, last_attempt)
	SELECT $3, subq.en, $4, $5, $6
	FROM subq
	ON CONFLICT (username, word_en) DO NOTHING`

	mock.ExpectExec(expectQuery).
		WithArgs(words[0], words[1], username, rating, attempts, lastAttempt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	sqlxDB := sqlx.NewDb(db, "pgx")
	client := &dbx.DBClient{
		DB: sqlxDB,
	}
	repo := New(client)

	_, repoErr := repo.AddWordToJar(ctx, username, rating, words...)

	assert.NoError(t, repoErr)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
