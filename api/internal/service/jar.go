package service

import (
	"context"
	"fmt"

	"github.com/maisiq/go-words-jar/internal/models"
)

const (
	StatusWordNew          string = "new"    // complete new word
	StatusWordWantToRepeat string = "medium" // already heard about this word
	StatusWordWellKnown    string = "wellknown"
)

var validStatuses = map[string]bool{
	StatusWordNew:          true,
	StatusWordWantToRepeat: true,
	StatusWordWellKnown:    true,
}

func IsValidStatus(status string) bool {
	_, ok := validStatuses[status]
	return ok
}

func (s *Service) GetUserWords(ctx context.Context, username string, params QueryParams) ([]models.Word, error) {
	words, err := s.repo.GetJarWords(ctx, username)
	if err != nil {
		s.log.Errorw("failed to get user words", "error", err)
		return nil, err
	}
	return words, nil
}

func (s *Service) AddWordsToJar(ctx context.Context, username, withStatus string, words ...string) (int64, error) {
	var rating float32

	if len(words) == 0 {
		return 0, fmt.Errorf("no words present")
	}

	switch withStatus {
	case StatusWordNew:
		rating = 0
	case StatusWordWantToRepeat:
		rating = 2.5
	case StatusWordWellKnown:
		rating = 5
	default:
		rating = 0
	}

	q, err := s.repo.AddWordToJar(ctx, username, rating, words...) // return words added
	if err != nil {
		s.log.Errorw("failed to add words to jar", "error", err, "words", words)
		return 0, err
	}
	return q, nil
}
