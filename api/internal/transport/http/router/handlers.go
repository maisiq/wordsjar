package router

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	UserInfo(c *gin.Context)
	CreateUser(c *gin.Context)
	Authenticate(c *gin.Context)
}

type ExamHandler interface {
	TestUser(c *gin.Context)
	VerifyWord(c *gin.Context)
}

type JarHandler interface {
	GetUserWords(c *gin.Context)
	AddWordToJar(c *gin.Context)
}

type WordsHandler interface {
	WordList(c *gin.Context)
	Word(c *gin.Context)
	AddWord(c *gin.Context)
}

type Handlers struct {
	Auth  AuthHandler
	Exam  ExamHandler
	Jar   JarHandler
	Words WordsHandler
}

type Middlewares struct {
	IsAdmin gin.HandlerFunc
	Auth    gin.HandlerFunc
	CORS    gin.HandlerFunc
}
