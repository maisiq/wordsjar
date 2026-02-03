package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maisiq/go-words-jar/internal/transport/http/middleware"
)

type UserInfoResponse struct {
	Username string `json:"username"`
}

func (h *AuthHandler) UserInfo(c *gin.Context) {
	username, ok := middleware.GetUsernameFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"details": "internal error"})
		return
	}

	info, err := h.service.UserInfo(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "internal error"})
		return
	}

	c.JSON(http.StatusOK, UserInfoResponse{
		Username: info.Username,
	})
}
