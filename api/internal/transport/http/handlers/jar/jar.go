package jar

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maisiq/go-words-jar/internal/service"
	"github.com/maisiq/go-words-jar/internal/transport/http/handlers"
	"github.com/maisiq/go-words-jar/internal/transport/http/middleware"
)

func (h *JarHandler) GetUserWords(c *gin.Context) {
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
	words, err := h.service.GetUserWords(c.Request.Context(), username, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, words)
}

type AddWordToJarRequest struct {
	WordEN string `form:"word_en" json:"word_en" binding:"required"`
	Status string `form:"status" json:"status"`
}

func (h *JarHandler) AddWordToJar(c *gin.Context) {
	username, ok := middleware.GetUsernameFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"details": "internal error"})
		return
	}

	var data AddWordToJarRequest
	err := c.ShouldBind(&data)

	if err != nil {
		payload := handlers.PayloadErrorJSON(err)
		c.JSON(http.StatusBadRequest, payload)
		return
	}

	//TODO: hold buffer overflow if too large text
	//TODO: split words by space and comma; hence, user can pass text and add all words from it
	r := strings.NewReplacer("\r\n", " ", "\n", " ", ",", "") //TODO: hold case: "word,word"
	rawWords := r.Replace(data.WordEN)
	words := strings.Split(rawWords, " ")

	var status string

	if service.IsValidStatus(data.Status) {
		status = data.Status
	} else {
		status = service.StatusWordNew
	}

	count, err := h.service.AddWordsToJar(c.Request.Context(), username, status, words...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"validation_error": false,
			"detail":           "internal error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"words_count": count,
	})
}
