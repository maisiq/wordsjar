package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/maisiq/go-words-jar/internal/service"
	"github.com/maisiq/go-words-jar/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestService_AddWordsToJar(t *testing.T) {
	log := zap.NewNop().Sugar()
	ctx := context.Background()

	username := "testuser"

	mc := minimock.NewController(t)
	repoMock := mocks.NewRepositoryMock(mc)
	s := service.New(log, repoMock)

	cases := []struct {
		name       string
		status     string
		words      []string
		wordsCount int64
		err        error

		expectRepoCall     bool
		repoCallWithWords  []string
		repoCallWithRating float32
		repoReturnsValue   int64
		repoReturnsErr     error
	}{
		{
			"success", service.StatusWordNew, []string{"word1", "word2"}, 2, nil,
			true, []string{"word1", "word2"}, 0, 2, nil,
		},
		{
			"empty status sets word as new(rating 0)", "", []string{"word1", "word2"}, 2, nil,
			true, []string{"word1", "word2"}, 0, 2, nil,
		},
		{
			"no words returns err", "", []string{""}, 0, fmt.Errorf("error"), // check error
			false, nil, 0, 0, nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.expectRepoCall {
				repoMock.AddWordToJarMock.Expect(ctx, username, c.repoCallWithRating, c.repoCallWithWords...).
					Return(c.repoReturnsValue, c.repoReturnsErr)
			}

			count, err := s.AddWordsToJar(ctx, username, c.status, c.words...)

			if c.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, c.wordsCount, count)

		})
	}

}
