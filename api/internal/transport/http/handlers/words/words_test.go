package words

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	httpx "github.com/maisiq/go-words-jar/internal/transport/http"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers/words/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWordsHandler_AddWord(t *testing.T) {
	ms := mocks.NewServiceMock(t)
	h := NewWordsrHandler(ms)
	regularUser := "testuser"

	cases := []struct {
		name              string
		serviceCall       bool
		serviceReturnsErr error
		payload           any
		body              string
	}{
		{"invalid payload returns bad request", false, nil, `just a string`, `{"validation_error":true,"errors":{"en":"field is required"}}`},
		{"field ru as not array returns validation error", false, nil, map[string]any{"en": "word", "ru": "word", "transcription": "trans"}, `{"validation_error":true,"errors":{"ru":"field should be array of string"}}`},
		{"required fields are empty", false, nil, map[string]any{"en": "", "ru": []string{""}, "transcription": ""}, `{"validation_error":true,"errors":{"en":"field is required", "ru":"field is required", "transcription":"field is required"}}`},
		{"some fields are not present", false, nil, map[string]any{"en": "word", "ru": []string{"слово"}, "transcription": ""}, `{"validation_error":true,"errors":{"transcription":"field is required"}}`},
		{"no payload returns validation error", false, nil, nil, `{"validation_error":true,"errors":{"en":"field is required", "ru":"field is required", "transcription":"field is required"}}`},
		{"service returns error", true, fmt.Errorf("bad"), map[string]any{"en": "word", "ru": []string{"слово"}, "transcription": "trans"}, `{"validation_error":false,"detail":"internal error"}`},
		{"add word with success", true, nil, map[string]any{"en": "word", "ru": []string{"слово"}, "transcription": "trans"}, `{}`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.serviceCall {
				ms.AddWordMock.Return(c.serviceReturnsErr)
			}
			buf := new(bytes.Buffer)
			_ = json.NewEncoder(buf).Encode(c.payload)

			rr := httptest.NewRecorder()
			tctx, _ := gin.CreateTestContext(rr)
			tctx.Set(httpx.UsernameContextKey, regularUser)

			request, _ := http.NewRequest("POST", "/words", buf)
			request.Header.Add("Content-Type", "application/json")
			tctx.Request = request

			h.AddWord(tctx)

			dataBytes, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %s", err)
			}

			assert.JSONEq(t, c.body, string(dataBytes))
		})
	}
}
