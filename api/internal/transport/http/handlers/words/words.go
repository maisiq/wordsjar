package words

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	errx "github.com/maisiq/go-words-jar/internal/errors"
	"github.com/maisiq/go-words-jar/internal/service"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers"
)

func (h *WordsHandler) WordList(c *gin.Context) {
	perPageStr := c.DefaultQuery("per_page", "10")
	lastID := c.DefaultQuery("last_id", "")
	prev := c.DefaultQuery("prev", "")

	perPageInt, err := strconv.Atoi(perPageStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "per_page query should be uint8"})
	}

	params := service.QueryParams{
		Limit: uint8(perPageInt),
	}

	if lastID != "" {
		next := true
		if prev != "" {
			next = false
		}
		params.Pagination = &service.Pagination{
			Next:    next,
			Pointer: lastID,
		}
	}

	words, err := service.Paginate(c.Request.Context(), params, h.service.WordList)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, words)
}

func (h *WordsHandler) Word(c *gin.Context) {
	name, ok := c.Params.Get("word")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "invalid request"})
		return
	}

	word, err := h.service.GetWordByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, word)
}

type AddWordRequest struct {
	EN            string   `json:"en" form:"en" binding:"required"`
	RU            []string `json:"ru" form:"ru" binding:"required,min=1,dive,required"`
	Transcription string   `json:"transcription" form:"transcription" binding:"required"`
}

func (h *WordsHandler) AddWord(c *gin.Context) {
	var data AddWordRequest
	err := c.ShouldBind(&data)

	if err != nil {
		payload := handlers.PayloadErrorJSON(err)
		c.JSON(http.StatusBadRequest, payload)
		return
	}

	err = h.service.AddWord(c.Request.Context(), data.EN, data.RU, data.Transcription)
	if err != nil {
		switch {
		case errors.Is(err, errx.ErrWordAlreadyExists):
			c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error(), "validation_error": false})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"detail": "internal error", "validation_error": false})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
