package exam

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/maisiq/go-words-jar/internal/service"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers"
	"github.com/maisiq/go-words-jar/internal/transport/http/middleware"
)

func (h *ExamHandler) TestUser(c *gin.Context) {
	username, ok := middleware.GetUsernameFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"details": "internal error"})
		return
	}

	perPage := c.DefaultQuery("per_page", "12")
	limit, err := strconv.Atoi(perPage)
	if err != nil {
		limit = 12
	}

	params := service.QueryParams{
		Limit: uint8(limit),
	}

	words, err := h.service.GetTestUserWords(c.Request.Context(), username, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": "internal error"})
		return
	}

	c.JSON(http.StatusOK, words)
}

type VerifyWordRequest struct {
	Reverse bool   `form:"reverse" json:"reverse"`
	WordID  string `form:"word_id" json:"word_id" binding:"required"`
	Answer  string `form:"answer" json:"answer" binding:"required"`
}

func (h *ExamHandler) VerifyWord(c *gin.Context) {
	username, ok := middleware.GetUsernameFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"details": "internal error"})
		return
	}

	var data VerifyWordRequest
	err := c.ShouldBindJSON(&data)

	if err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			errorMsgs := handlers.GetValidationErrorMessages(verr)
			c.JSON(http.StatusBadRequest, gin.H{
				"validation_error": true,
				"errors":           errorMsgs,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"validation_error": false,
				"detail":           "bad request",
			})
		}
		return
	}

	passed, err := h.service.VerifyWord(c.Request.Context(), username, data.WordID, data.Answer, data.Reverse)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"validation_error": false,
			"detail":           "internal error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"passed": passed,
	})

}
