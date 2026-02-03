package exam

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
	"github.com/maisiq/go-words-jar/internal/service"
	httpx "github.com/maisiq/go-words-jar/internal/transport/http"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/exam/mocks"
	"github.com/stretchr/testify/assert"
)

func TestExamHandler_TestUser(t *testing.T) {
	ms := mocks.NewServiceMock(t)

	h := NewExamHandler(ms)
	username := "test_user"
	pw := []service.TestWord{}

	cases := []struct {
		name                  string
		queryParams           map[string]string
		expectedServiceParams service.QueryParams
		serviceReturnsWords   []service.TestWord
		serviceReturnsError   error
		expectedStatusCode    int
	}{
		{"per page 10 matches", map[string]string{"per_page": "10"}, service.QueryParams{Limit: 10}, pw, nil, http.StatusOK},
		{"per page 15 matches", map[string]string{"per_page": "15"}, service.QueryParams{Limit: 15}, pw, nil, http.StatusOK},
		{"per page 1000 trims", map[string]string{"per_page": "1000"}, service.QueryParams{Limit: 232}, pw, nil, http.StatusOK},
		{"per page with not a number returns default", map[string]string{"per_page": "invalid"}, service.QueryParams{Limit: 12}, pw, nil, http.StatusOK},
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

		ms.GetTestUserWordsMock.ExpectParamsParam3(c.expectedServiceParams).Return(c.serviceReturnsWords, c.serviceReturnsError)

		h.TestUser(tctx)
		assert.Equal(t, c.expectedStatusCode, rr.Result().StatusCode)
	}
}

func TestExamHandler_VerifyWord(t *testing.T) {
	ms := mocks.NewServiceMock(t)
	h := NewExamHandler(ms)
	username := "test_user"

	cases := []struct {
		name                 string
		serviceCall          bool
		serviceRerurnsResult bool
		serviceReturnsErr    error
		payload              any
		body                 string
	}{
		{"required field is empty", false, false, nil, map[string]string{"answer": "", "word_id": ""}, `{"validation_error":true,"errors":{"answer":"field is required", "word_id":"field is required"}}`},
		{"no payload returns validation error", false, false, nil, nil, `{"validation_error":true,"errors":{"answer":"field is required", "word_id":"field is required"}}`},
		{"service returns unpredictable error", true, false, fmt.Errorf("bad"), map[string]string{"answer": "слово", "word_id": "some-id"}, `{"validation_error":false,"detail":"internal error"}`},
		{"verify word with success", true, true, nil, map[string]string{"answer": "слово", "word_id": "some-id"}, `{"passed":true}`},
		{"verify word with failure", true, false, nil, map[string]string{"answer": "слово", "word_id": "some-id"}, `{"passed":false}`},
		{"invalid payload returns bad request", false, false, nil, `just a string`, `{"validation_error":false,"detail":"bad request"}`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.serviceCall {
				ms.VerifyWordMock.Return(c.serviceRerurnsResult, c.serviceReturnsErr)
			}
			buf := new(bytes.Buffer)
			_ = json.NewEncoder(buf).Encode(c.payload)

			rr := httptest.NewRecorder()
			tctx, _ := gin.CreateTestContext(rr)
			tctx.Set(httpx.UsernameContextKey, username)

			request, _ := http.NewRequest("POST", "/test", buf)
			request.Header.Add("Content-Type", "application/json")
			tctx.Request = request
			h.VerifyWord(tctx)

			dataBytes, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %s", err)
			}

			assert.JSONEq(t, c.body, string(dataBytes))
		})
	}
}
