package auth_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maisiq/go-words-jar/internal/config"
	authService "github.com/maisiq/go-words-jar/internal/service/auth"
	httpx "github.com/maisiq/go-words-jar/internal/transport/http"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/auth"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/auth/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_UserInfo(t *testing.T) {
	service := mocks.NewServiceMock(t)
	regularUser := "testuser"
	h := auth.NewAuthHandler(config.HTTP{}, service)

	cases := []struct {
		name               string
		username           string
		serviceReturnsData authService.UserInfo
		serviceReturnsErr  error
		responseCode       int
		responseData       string
	}{
		{"success", "user", authService.UserInfo{Username: "user"}, nil, http.StatusOK, "{\"username\":\"user\"}"},
		{"user doesn't exist", "user1", authService.UserInfo{}, fmt.Errorf("no found"), http.StatusInternalServerError, "{\"detail\":\"internal error\"}"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			service.UserInfoMock.Return(c.serviceReturnsData, c.serviceReturnsErr)

			rr := httptest.NewRecorder()
			tctx, _ := gin.CreateTestContext(rr)
			tctx.Set(httpx.UsernameContextKey, regularUser)

			request, _ := http.NewRequest("GET", "/user", nil)

			request.Header.Add("Content-Type", "application/json")
			tctx.Request = request

			h.UserInfo(tctx)

			assert.Equal(t, rr.Result().StatusCode, c.responseCode)

			dataBytes, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Logf("failed to read body: %s", err)
			}

			assert.Equal(t, c.responseData, string(dataBytes))
		})

	}

}
