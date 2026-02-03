package jar_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maisiq/go-words-jar/internal/models"
	"github.com/maisiq/go-words-jar/internal/service"
	httpx "github.com/maisiq/go-words-jar/internal/transport/http"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/jar"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/jar/mocks"
	"github.com/stretchr/testify/assert"
)

func TestJarHandler_AddWordToJar(t *testing.T) {
	service := mocks.NewServiceMock(t)
	h := jar.NewJarHandler(service)

	username := "test_user"

	cases := []struct {
		name                string
		serviceCall         bool
		serviceReturnsErr   error
		serviceReturnsValue int64
		payload             any
		body                string
	}{
		{"required field is empty", false, nil, 0, map[string]string{"word_en": ""}, `{"validation_error":true,"errors":{"word_en":"field is required"}}`},
		{"no payload returns validation error", false, nil, 0, nil, `{"validation_error":true,"errors":{"word_en":"field is required"}}`},
		{"service returns error", true, fmt.Errorf("bad"), 0, map[string]string{"word_en": "add"}, `{"validation_error":false,"detail":"internal error"}`},
		{"add word with success", true, nil, 1, map[string]string{"word_en": "add"}, `{"words_count":1}`},
		{"add words with status", true, nil, 1, map[string]string{"word_en": "add", "status": "wellknown"}, `{"words_count":1}`},
		{"invalid payload returns validation error with first field", false, nil, 0, `just a string`, `{"validation_error":true,"errors":{"word_en":"field is required"}}`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.serviceCall {
				service.AddWordsToJarMock.Return(c.serviceReturnsValue, c.serviceReturnsErr)
			}
			buf := new(bytes.Buffer)
			_ = json.NewEncoder(buf).Encode(c.payload)

			rr := httptest.NewRecorder()
			tctx, _ := gin.CreateTestContext(rr)
			tctx.Set(httpx.UsernameContextKey, username)

			request, _ := http.NewRequest("POST", "/words", buf)
			request.Header.Add("Content-Type", "application/json")
			tctx.Request = request

			h.AddWordToJar(tctx)

			dataBytes, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %s", err)
			}

			assert.JSONEq(t, c.body, string(dataBytes))
		})
	}
}

func TestJarHandler_GetUserWords(t *testing.T) {
	ms := mocks.NewServiceMock(t)
	h := jar.NewJarHandler(ms)

	username := "test_user"

	pw := []models.Word{}

	cases := []struct {
		name                  string
		queryParams           map[string]string
		expectedServiceParams service.QueryParams
		serviceReturnsWords   []models.Word
		serviceReturnsError   error
		expectedStatusCode    int
	}{
		{"per page 10 matches", map[string]string{"per_page": "10"}, service.QueryParams{Limit: 10}, pw, nil, http.StatusOK},
		{"per page 15 matches", map[string]string{"per_page": "15"}, service.QueryParams{Limit: 15}, pw, nil, http.StatusOK},
		{"per page 1000 trims", map[string]string{"per_page": "1000"}, service.QueryParams{Limit: 232}, pw, nil, http.StatusOK},
	}

	for _, c := range cases {
		q := url.Values{}

		rr := httptest.NewRecorder()
		tctx, _ := gin.CreateTestContext(rr)
		tctx.Set(httpx.UsernameContextKey, username)

		request, _ := http.NewRequest("GET", "/test", nil)
		request.Header.Add("Content-Type", "application/json")

		for k, v := range c.queryParams {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
		tctx.Request = request

		ms.GetUserWordsMock.ExpectParamsParam3(c.expectedServiceParams).Return(c.serviceReturnsWords, c.serviceReturnsError)

		h.GetUserWords(tctx)

		assert.Equal(t, c.expectedStatusCode, rr.Result().StatusCode)
	}

}
