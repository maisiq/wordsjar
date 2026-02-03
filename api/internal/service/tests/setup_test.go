package service_test

import (
	"slices"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/maisiq/go-words-jar/internal/models"
)

var TestWordsASC []models.Word
var TestWordsDESC []models.Word

func TestMain(m *testing.M) {
	maxWords := 20

	wordsAsc := make([]models.Word, 0, maxWords)
	for i := 0; i < maxWords; i++ {
		wordsAsc = append(wordsAsc, models.Word{
			ID:            gofakeit.UUID(),
			EN:            gofakeit.Word(),
			RU:            []string{gofakeit.Word(), gofakeit.Word()},
			Transcription: gofakeit.Word(),
		})
	}

	wordsDesc := make([]models.Word, maxWords)

	copy(wordsDesc, wordsAsc)

	slices.SortFunc(wordsDesc, func(a models.Word, b models.Word) int {
		if a.ID <= b.ID {
			return 1
		} else {
			return -1
		}
	})

	slices.SortFunc(wordsAsc, func(a models.Word, b models.Word) int {
		if a.ID >= b.ID {
			return 1
		} else {
			return -1
		}
	})

	TestWordsASC = wordsAsc
	TestWordsDESC = wordsDesc
}
