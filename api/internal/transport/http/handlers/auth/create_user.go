package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maisiq/go-words-jar/internal/errors"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var data CreateUserRequest
	err := c.ShouldBind(&data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid request"})
		return
	}

	err = h.service.CreateUser(c.Request.Context(), data.Username, data.Password)
	if err != nil {
		switch err {
		case errors.ErrUserAlreadyExists:
			c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
