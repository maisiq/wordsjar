package repository

import (
	"context"

	"github.com/maisiq/go-words-jar/internal/db"
	"github.com/maisiq/go-words-jar/internal/models"
	"github.com/maisiq/go-words-jar/internal/service"
)

var _ IRepository = (*Repository)(nil)

type IRepository interface {
	GetUserWords(ctx context.Context, username string, filters ...service.Filter) ([]models.Word, error)
	GetWordByID(ctx context.Context, id string) (models.Word, error)
	GetWordByName(ctx context.Context, wordName string) (models.Word, error)
	WordList(context.Context, *service.QueryParams) ([]models.Word, error)
	AddWord(context.Context, models.Word) error

	GetUserWord(ctx context.Context, wordID string, username string) (*models.UserWord, error)
	AddWordToJar(ctx context.Context, username string, setRating float32, words ...string) (wordCount int64, err error)
	UpdateUserWord(ctx context.Context, userWord *models.UserWord) error
}

type Repository struct {
	client *db.DBClient
}

func New(c *db.DBClient) *Repository {
	return &Repository{
		client: c,
	}
}
