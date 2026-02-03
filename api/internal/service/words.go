package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	errx "github.com/maisiq/go-words-jar/internal/errors"
	"github.com/maisiq/go-words-jar/internal/models"
)

func (s *Service) GetWordByName(ctx context.Context, word string) (models.Word, error) {
	w, err := s.repo.GetWordByName(ctx, word)
	if err != nil {
		return models.Word{}, err
	}
	return w, nil
}

func (s *Service) AddWord(ctx context.Context, en string, ru []string, transcription string) error {
	id := uuid.New()
	translations := make([]string, 0, len(ru))

	for _, trans := range ru {
		val := strings.Trim(trans, " ")
		if val != "" {
			translations = append(translations, val)
		}
	}

	if len(translations) == 0 {
		return fmt.Errorf("ru field should contains atleast one element")
	}

	word := models.Word{
		ID:            id.String(),
		EN:            en,
		RU:            ru,
		Transcription: transcription,
	}

	repoErr := s.repo.AddWord(ctx, word)
	if repoErr != nil {
		switch {
		case errors.Is(repoErr, errx.ErrWordAlreadyExists):
			return repoErr
		default:
			s.log.Errorw("failed to add word", "detail", repoErr)
			return errx.ErrInternal
		}
	}
	return nil
}

func (s *Service) WordList(ctx context.Context, params QueryParams) ([]models.Word, error) {
	if params.SortBy == "" {
		params.SortBy = FieldID
	}
	if params.Limit == 0 {
		params.Limit = 20
	}

	words, err := s.repo.WordList(ctx, &params)
	if err != nil {
		s.log.Errorw("failed to get word list", "error", err)
		return nil, err
	}

	return words, nil
}
