package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	notAllowedHost := "not-allowed.host"

	AllowedHosts = []string{
		"allowed-host.ru",
		"one-more-allowed.com",
	}

	cases := []struct {
		name                string
		method              string
		withHeaders         []string
		withHeaderValues    []string
		checkHeaderName     string
		expectedHeaderValue string
	}{
		{"not allowed host", "GET", []string{"Origin"}, []string{notAllowedHost}, "Access-Control-Allow-Origin", ""},
		{"only get, options methods allowed when host not in allowed host list", "GET", []string{"Origin"}, []string{notAllowedHost}, "Access-Control-Allow-Methods", "GET, OPTIONS"},
		{"allow various methods when host in allowed host list", "GET", []string{"Origin"}, []string{AllowedHosts[0]}, "Access-Control-Allow-Methods", "GET, POST, PUT, PUTCH, DELETE, OPTIONS"},
		{"allowed host return allow origin header with right value", "GET", []string{"Origin"}, []string{AllowedHosts[0]}, "Access-Control-Allow-Origin", AllowedHosts[0]},
		{"allowed host with credenetials", "GET", []string{"Origin", "Access-Control-Allow-Headers"}, []string{AllowedHosts[0], "Authorization"}, "Access-Control-Allow-Credentials", "true"},
	}

	r := gin.New()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			ctx := gin.CreateTestContextOnly(response, r)
			ctx.Request = httptest.NewRequest(c.method, "/", nil)

			if len(c.withHeaders) != len(c.withHeaderValues) {
				t.Errorf("len set headers and their values not match")
			}

			for i := 0; i < len(c.withHeaders); i++ {
				ctx.Request.Header.Set(c.withHeaders[i], c.withHeaderValues[i])
			}

			CORSMiddleware(ctx)

			actualHeaderValue := response.Header().Get(c.checkHeaderName)
			assert.Equal(t, c.expectedHeaderValue, actualHeaderValue)
		})
	}

}
