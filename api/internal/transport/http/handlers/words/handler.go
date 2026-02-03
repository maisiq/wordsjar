package words

import (
	"context"

	"github.com/maisiq/go-words-jar/internal/models"
	"github.com/maisiq/go-words-jar/internal/service"
)

type Service interface {
	GetWordByName(context.Context, string) (models.Word, error)
	AddWord(ctx context.Context, en string, ru []string, transcription string) error
	WordList(context.Context, service.QueryParams) ([]models.Word, error)
}

type WordsHandler struct {
	service Service
}

func NewWordsrHandler(service Service) *WordsHandler {
	return &WordsHandler{
		service: service,
	}
}
