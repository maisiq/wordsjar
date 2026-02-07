package repository

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/maisiq/go-words-jar/internal/models"
	"github.com/maisiq/go-words-jar/internal/service"
)

func (repo *Repository) GetUserWords(ctx context.Context, username string, filters ...service.Filter) ([]models.Word, error) {
	result := []models.Word{}
	queryFilters := &service.UserWordsFilter{}

	words := sq.Select("w.id", "w.en", "w.ru", "w.transcription").
		From("words w").
		InnerJoin("user_words uw ON uw.word_en = w.en").
		InnerJoin("users u ON u.username = uw.username").
		Where("u.username = $1", username).
		PlaceholderFormat(sq.Dollar)

	for _, filter := range filters {
		filter(queryFilters)
	}

	if queryFilters.TestMode {
		words = words.Where("last_attempt + INTERVAL '1 day' < now() AT TIME ZONE 'UTC'").
			Where("knowledge_rating < $2", 5.0).
			OrderBy("uw.knowledge_rating ASC")
	} else {
		words = words.OrderBy("id")
	}

	if queryFilters.WordID != "" {
		ph := "$2"
		if queryFilters.TestMode {
			ph = "$3"
		}
		words = words.Where(fmt.Sprintf("w.id = %s", ph), queryFilters.WordID)
	}

	query, args, err := words.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}

	rows, err := repo.client.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from db: %w", err)
	}

	for rows.Next() {
		var word Word
		err := rows.StructScan(&word)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row into Word: %w", err)
		}

		w := models.Word{
			ID:            word.ID,
			EN:            word.EN,
			RU:            word.RuTranslations(),
			Transcription: word.Transcription,
		}

		result = append(result, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *Repository) AddWordToJar(ctx context.Context, username string, setRating float32, wordEN ...string) (int64, error) {
	// prevents add words which don't exist in db without error
	baseQuery := `WITH subq as NOT MATERIALIZED (
		SELECT en
		FROM words
		WHERE en IN (?)
	)
	INSERT INTO user_words (username, word_en, knowledge_rating, consecutive_success_attempts, last_attempt)
	SELECT ?, subq.en, ?, ?, ?
	FROM subq
	ON CONFLICT (username, word_en) DO NOTHING`

	// ensure that user can immediately get words in his tests
	lastAttempt := time.Now().UTC().Add(-time.Hour * 24)

	newQ, args, err := sqlx.In(baseQuery, wordEN, username, setRating, 1, lastAttempt)
	if err != nil {
		return 0, err
	}
	newQ = sqlx.Rebind(sqlx.DOLLAR, newQ)

	result, err := repo.client.DB.ExecContext(ctx, newQ, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to add word: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to add word: %w", err) // change & test
	}

	return count, nil
}

// wordID or wordEN?? worid can help with other langs
func (repo *Repository) GetUserWord(ctx context.Context, wordID string, username string) (*models.UserWord, error) {
	query := `SELECT w.id as word_id, uw.username, uw.knowledge_rating, uw.consecutive_success_attempts
FROM words w
INNER JOIN user_words uw ON uw.word_en = w.en
INNER JOIN users u ON u.username = uw.username
WHERE u.username = $1 AND w.id = $2`

	row := repo.client.DB.QueryRowxContext(ctx, query, username, wordID)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed to get rows from db: %w", row.Err())
	}

	var userWord UserWord

	err := row.StructScan(&userWord)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row into UserWord: %w", err)
	}

	word := &models.UserWord{
		WordID:      userWord.WordID,
		Username:    userWord.Username,
		Rating:      userWord.Rating,
		Attempts:    userWord.Attempts,
		LastAttempt: userWord.LastAttempt,
	}
	return word, nil
}

func (repo *Repository) UpdateUserWord(ctx context.Context, userWord *models.UserWord) error {
	query := `UPDATE user_words uw
SET knowledge_rating = :rating, consecutive_success_attempts = :attempts, last_attempt = :last_attempt
FROM words w
WHERE uw.word_en = w.en AND uw.username = :username AND w.id = :word_id`

	_, err := repo.client.DB.NamedExecContext(ctx, query, map[string]interface{}{
		"rating":       userWord.Rating,
		"attempts":     userWord.Attempts,
		"username":     userWord.Username,
		"word_id":      userWord.WordID,
		"last_attempt": userWord.LastAttempt,
	})

	if err != nil {
		return fmt.Errorf("failed to update UserWord: %w", err)
	}

	return nil
}
