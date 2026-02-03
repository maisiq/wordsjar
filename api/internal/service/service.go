package service

import (
	"context"

	"github.com/maisiq/go-words-jar/internal/logger"
	"github.com/maisiq/go-words-jar/internal/models"
)

type JarRepository interface {
	GetUserWords(ctx context.Context, username string, filters ...Filter) ([]models.Word, error)
	AddWordToJar(ctx context.Context, username string, setRating float32, words ...string) (int64, error)
}

type WordsRepository interface {
	GetWordByName(ctx context.Context, wordName string) (models.Word, error)
	AddWord(context.Context, models.Word) error
	WordList(context.Context, *QueryParams) ([]models.Word, error)
}

type ExamRepository interface {
	GetUserWords(ctx context.Context, username string, filters ...Filter) ([]models.Word, error)
	GetUserWord(ctx context.Context, wordID string, username string) (*models.UserWord, error)
	GetWordByID(ctx context.Context, id string) (models.Word, error)
	UpdateUserWord(ctx context.Context, userWord *models.UserWord) error
}

type Repository interface {
	JarRepository
	WordsRepository
	ExamRepository
}

type Service struct {
	repo Repository
	log  logger.Logger
}

func New(logger logger.Logger, repo Repository) *Service {
	return &Service{
		repo: repo,
		log:  logger,
	}
}
