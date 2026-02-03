package auth

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/maisiq/go-words-jar/internal/service/auth/mocks"
	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
)

func TestLogout(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop().Sugar()

	storageMock := mocks.NewTokenStorageMock(t)
	service := NewAuthService(log, nil, storageMock, nil)
	ttl := int64(10)

	t.Run("service can save tokens hence invalidates it", func(t *testing.T) {
		tokens := Tokens{
			Access:  "access token value",
			Refresh: "refresh token value",
		}

		storageMock.StoreMock.When(minimock.AnyContext, tokens.Access, "true", ttl).Then(nil)
		storageMock.StoreMock.When(minimock.AnyContext, tokens.Refresh, "true", ttl).Then(nil)

		err := service.Logout(ctx, tokens)

		assert.NoError(t, err)
	})
}
