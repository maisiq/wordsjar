package jar

import (
	"context"

	"github.com/maisiq/go-words-jar/internal/models"
	"github.com/maisiq/go-words-jar/internal/service"
)

type Service interface {
	GetUserWords(ctx context.Context, username string, params service.QueryParams) ([]models.Word, error)
	AddWordsToJar(ctx context.Context, username string, withStatus string, words ...string) (int64, error)
}

type JarHandler struct {
	service Service
}

func NewJarHandler(uService Service) *JarHandler {
	return &JarHandler{
		service: uService,
	}
}
