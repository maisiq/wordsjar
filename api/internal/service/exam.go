package service

import (
	"context"
	"strings"
	"time"

	"github.com/maisiq/go-words-jar/internal/models"
)

const (
	successPoints float32 = 0.5
	failurePoints float32 = 0.25
	maxRating     float32 = 5.0
	minRating     float32 = 0
)

type TestWord struct {
	models.Word
	Reverse bool `json:"reverse"`
}

func (s *Service) GetTestUserWords(ctx context.Context, username string, params QueryParams) ([]TestWord, error) {
	words, err := s.repo.GetJarWords(ctx, username, WithTestMode())
	if err != nil {
		s.log.Errorw("failed to get paginated words", "error", err)
		return nil, err
	}

	testWords := make([]TestWord, 0, len(words))

	for idx, word := range words {
		testWord := TestWord{
			Word: word,
		}
		if idx%2 == 0 {
			testWord.Reverse = true
		}
		testWords = append(testWords, testWord)
	}

	return testWords, nil
}

func (s *Service) VerifyWord(ctx context.Context, username, wordID, answer string, reverse bool) (bool, error) {
	word, err := s.repo.GetWordByID(ctx, wordID)
	if err != nil {
		return false, err
	}

	answer = strings.ToLower(answer)
	answer = strings.Trim(answer, " ")

	var result bool

	if reverse {
		if word.EN == answer {
			result = true
		}
	} else {
		for _, w := range word.RU {
			if w == answer {
				result = true
				break
			}
		}
	}

	userWord, err := s.repo.GetUserWord(ctx, word.ID, username)
	if err != nil {
		s.log.Errorw("failed to get word rating for the user", "error", err.Error())
		// even if there is an error user could be right anyway
		return result, err
	}

	if result {
		newRating := userWord.Rating + successPoints*float32(userWord.Attempts)
		if newRating > maxRating {
			userWord.Rating = maxRating
		} else {
			userWord.Rating = newRating
		}
		userWord.Attempts++
	} else {
		newRating := userWord.Rating - failurePoints
		if newRating < minRating {
			userWord.Rating = minRating
		} else {
			userWord.Rating = newRating
		}
		userWord.Attempts = 1
	}

	userWord.LastAttempt = time.Now().UTC()

	err = s.repo.UpdateUserWord(ctx, userWord)
	if err != nil {
		s.log.Errorw("failed to set word rating with attempts", "error", err.Error())
		// even if there is an error user could be right anyway
		return result, err
	}
	return result, nil
}
