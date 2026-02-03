package exam

import (
	"context"

	"github.com/maisiq/go-words-jar/internal/service"
)

type Service interface {
	GetTestUserWords(ctx context.Context, username string, params service.QueryParams) ([]service.TestWord, error)
	VerifyWord(ctx context.Context, username, wordID, answer string, enCheck bool) (bool, error)
}

type ExamHandler struct {
	service Service
}

func NewExamHandler(service Service) *ExamHandler {
	return &ExamHandler{
		service: service,
	}
}
