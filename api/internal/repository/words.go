package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/maisiq/go-words-jar/internal/models"
	"github.com/maisiq/go-words-jar/internal/service"
)

type field string

const (
	FieldID field = "id"
	FieldEN field = "en"
)

type Pagination struct {
	Pointer interface{}
	Next    bool
}

type QueryParams struct {
	Limit      uint8
	SortBy     field
	Desc       bool
	Pagination *Pagination
}

func (repo *Repository) WordList(ctx context.Context, params *service.QueryParams) ([]models.Word, error) {
	var words []models.Word

	builder := sq.Select("id", "en", "ru", "transcription").PlaceholderFormat(sq.Dollar).
		From("words")

	if params != nil {
		builder = builder.Limit(uint64(params.Limit))
		orderBy := string(params.SortBy)
		sortField := orderBy

		if params.Pagination != nil {
			var whereCond string
			var sign string

			if params.Pagination.Next {
				if params.Desc {
					orderBy += " DESC"
					sign = "<"
				} else {
					sign = ">"
				}
			} else {
				// reversing sort if previous page
				if params.Desc {
					orderBy += " ASC"
					sign = ">"
				} else {
					sign = "<"
					orderBy += " DESC"
				}
			}
			if params.SortBy != service.FieldID {
				whereCond = fmt.Sprintf("%s %s (SELECT %s FROM words WHERE id = $1)", sortField, sign, sortField)
			} else {
				whereCond = fmt.Sprintf("id %s $1", sign)
			}
			builder = builder.Where(whereCond, params.Pagination.Pointer)
		} else {
			if params.Desc {
				orderBy += " DESC"
			}
		}
		builder = builder.OrderBy(orderBy)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}

	result, err := repo.client.DB.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		var word Word
		err := result.StructScan(&word)
		if err != nil {
			return nil, err
		}

		w := models.Word{
			ID:            word.ID,
			EN:            word.EN,
			RU:            word.RuTranslations(),
			Transcription: word.Transcription,
		}
		words = append(words, w)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func (repo *Repository) GetWordByID(ctx context.Context, id string) (models.Word, error) {
	query := `SELECT id, en, ru, transcription
FROM words
WHERE id = $1`

	row := repo.client.DB.QueryRowxContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return models.Word{}, err
	}
	var word Word
	err := row.StructScan(&word)
	if err != nil {
		return models.Word{}, err
	}

	return models.Word{
		ID:            word.ID,
		EN:            word.EN,
		RU:            word.RuTranslations(),
		Transcription: word.Transcription,
	}, nil
}

func (repo *Repository) GetWordByName(ctx context.Context, wordName string) (models.Word, error) {
	query := `SELECT id, en, ru, transcription
FROM words
WHERE en = $1`

	row := repo.client.DB.QueryRowxContext(ctx, query, wordName)
	if err := row.Err(); err != nil {
		return models.Word{}, err
	}
	var word Word
	err := row.StructScan(&word)
	if err != nil {
		return models.Word{}, err
	}

	return models.Word{
		ID:            word.ID,
		EN:            word.EN,
		RU:            word.RuTranslations(),
		Transcription: word.Transcription,
	}, nil
}
