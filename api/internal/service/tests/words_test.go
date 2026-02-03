package service_test

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/maisiq/go-words-jar/internal/models"
	"github.com/maisiq/go-words-jar/internal/service"
	"github.com/maisiq/go-words-jar/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestService_WordList(t *testing.T) {
	log := zap.NewNop().Sugar()
	ctx := context.Background()
	limit := 5
	mc := minimock.NewController(t)

	repoMock := mocks.NewRepositoryMock(mc)
	s := service.New(log, repoMock)

	// all cases should be in pagination
	t.Run("word list returns 1st page", func(t *testing.T) {
		repoMock.WordListMock.Return(TestWordsASC[:limit+1], nil)
		expected := service.Paginated[models.Word]{
			Items:      TestWordsASC[:limit],
			NextCursor: TestWordsASC[limit-1].ID,
			PrevCursor: "",
			HasNext:    true,
		}

		pw, err := s.WordList(ctx, service.QueryParams{
			Limit:  uint8(limit),
			SortBy: "", // not used
			Desc:   false,
		})

		assert.NoError(t, err, "got error from word list")
		assert.Equal(t, expected, pw)
	})

	t.Run("word list returns 2d page", func(t *testing.T) {
		repoMock.WordListMock.Return(TestWordsASC[limit:limit*2+1], nil)
		expected := service.Paginated[models.Word]{
			Items:      TestWordsASC[limit : 2*limit],
			NextCursor: TestWordsASC[2*limit-1].ID,
			PrevCursor: TestWordsASC[limit].ID,
			HasNext:    true,
		}

		pw, err := s.WordList(ctx, service.QueryParams{
			Limit:  uint8(limit),
			SortBy: "", // not used
			Desc:   false,
			Pagination: &service.Pagination{
				Next:    true,
				Pointer: TestWordsASC[limit-1].ID,
			},
		})

		assert.NoError(t, err, "got error from word list")
		assert.Equal(t, expected, pw)
	})

	t.Run("word list returns from page to previous page", func(t *testing.T) {
		var repoWords []models.Word
		expected := service.Paginated[models.Word]{
			Items:      TestWordsASC[:limit],
			NextCursor: TestWordsASC[limit-1].ID,
			PrevCursor: "",
			HasNext:    true,
		}

		pointer := TestWordsASC[limit].ID
		for _, w := range TestWordsDESC { // works only for 1st page
			if len(repoWords) == limit+1 {
				break
			}
			if w.ID < pointer {
				repoWords = append(repoWords, w)
			}
		}

		repoMock.WordListMock.Return(repoWords, nil)

		pw, err := s.WordList(ctx, service.QueryParams{
			Limit:  uint8(limit),
			SortBy: "", // not used
			Desc:   false,
			Pagination: &service.Pagination{
				Next:    false,
				Pointer: pointer,
			},
		})

		assert.NoError(t, err, "got error from word list")
		assert.Equal(t, expected, pw)
	})
}

func TestService_AddWord(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop().Sugar()

	word := models.Word{
		EN:            "word",
		RU:            []string{"слово"},
		Transcription: "trans",
	}
	repo := mocks.NewRepositoryMock(t)
	s := service.New(log, repo)

	t.Run("success", func(t *testing.T) {
		repo.AddWordMock.Return(nil)

		err := s.AddWord(ctx, word.EN, word.RU, word.Transcription)
		assert.NoError(t, err)
	})

	t.Run("error when RU slice has 0 len", func(t *testing.T) {
		err := s.AddWord(ctx, word.EN, []string{}, word.Transcription)

		assert.NotNil(t, err)
	})

	t.Run("error when RU slice has empty strings", func(t *testing.T) {
		err := s.AddWord(ctx, word.EN, []string{"", ""}, word.Transcription)

		assert.NotNil(t, err)
	})

}
