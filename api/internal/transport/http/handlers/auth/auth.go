package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authService "github.com/maisiq/go-words-jar/internal/service/auth"
	httpx "github.com/maisiq/go-words-jar/internal/transport/http"
)

type AutheticateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Authenticate(c *gin.Context) {
	var data AutheticateRequest
	err := c.ShouldBind(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid request"})
		return
	}

	tokens, err := h.service.Authenticate(c.Request.Context(), data.Username, data.Password)
	if err != nil {
		switch err {
		case authService.ErrAutheticationFailed:
			c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"detail": "internal error"})
		}
		return
	}

	c.SetCookie(httpx.AccessTokenKey, tokens.Access, h.cfg.Cookies.MaxAge, "/", h.cfg.Cookies.BaseDomain, h.cfg.Cookies.Secure, true)   // fix
	c.SetCookie(httpx.RefreshTokenKey, tokens.Refresh, h.cfg.Cookies.MaxAge, "/", h.cfg.Cookies.BaseDomain, h.cfg.Cookies.Secure, true) // fix

	c.JSON(http.StatusOK, tokens)
}
