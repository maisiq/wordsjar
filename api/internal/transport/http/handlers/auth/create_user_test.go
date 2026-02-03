package auth_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maisiq/go-words-jar/internal/config"
	"github.com/maisiq/go-words-jar/internal/errors"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/auth"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/auth/mocks"
	"github.com/maisiq/go-words-jar/internal/transport/http/router"
	routerMocks "github.com/maisiq/go-words-jar/internal/transport/http/router/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_CreateUser(t *testing.T) {
	service := mocks.NewServiceMock(t)
	routes := router.New(&router.Handlers{
		Auth:  auth.NewAuthHandler(config.HTTP{}, service),
		Words: routerMocks.NewWordsHandlerMock(t),
		Exam:  routerMocks.NewExamHandlerMock(t),
		Jar:   routerMocks.NewJarHandlerMock(t),
	}, &router.Middlewares{}, router.RouterParams{})
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	cases := []struct {
		name           string
		serviceCalled  bool
		serviceReturns error
		payload        map[string]string
		responseCode   int
		reponseData    string
	}{
		{"success", true, nil, map[string]string{"username": "test_user", "password": "mypassword"}, http.StatusOK, "{}"},
		{"user already exists", true, errors.ErrUserAlreadyExists, map[string]string{"username": "test_user", "password": "mypassword"}, http.StatusBadRequest, "{\"detail\":\"user already exists\"}"},
		{"internal error", true, errors.ErrInternal, map[string]string{"username": "test_user", "password": "mypassword"}, http.StatusInternalServerError, "{\"detail\":\"internal error\"}"},
		{"bad data request", false, nil, map[string]string{"username": "test_user"}, http.StatusBadRequest, "{\"detail\":\"invalid request\"}"},
		{"request with empty fields returns bad request", false, nil, map[string]string{"username": "", "password": ""}, http.StatusBadRequest, "{\"detail\":\"invalid request\"}"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.serviceCalled {
				service.CreateUserMock.Return(c.serviceReturns)
			}

			buf := new(bytes.Buffer)
			_ = json.NewEncoder(buf).Encode(c.payload)
			resp, err := ts.Client().Post(ts.URL+"/signup", "application/json", buf)
			if err != nil {
				t.Fatal("failed to request /singup")
			}
			defer resp.Body.Close()

			assert.Equal(t, resp.StatusCode, c.responseCode)

			dataBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Logf("failed to read body: %s", err)
			}

			assert.Equal(t, c.reponseData, string(dataBytes))
		})

	}

}
