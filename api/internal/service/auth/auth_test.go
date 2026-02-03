package auth

import (
	"testing"

	"github.com/maisiq/go-words-jar/internal/service/auth/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestUserService(t *testing.T) {
	log := zap.NewNop().Sugar()
	repo := mocks.NewRepositoryMock(t)
	service := NewAuthService(log, repo, nil, nil)

	t.Run("service can create user", func(t *testing.T) {
		username := "test_user"
		plainPassword := "test_password"
		repo.AddUserMock.Return(nil)

		err := service.CreateUser(t.Context(), username, plainPassword)

		assert.NoError(t, err)
	})
}
