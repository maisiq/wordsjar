package router

import (
	"github.com/gin-gonic/gin"
	"github.com/maisiq/go-words-jar/internal/transport/http/middleware"
)

type RouterParams struct {
	Debug bool
}

func New(handlers *Handlers, middlewares *Middlewares, params RouterParams) *gin.Engine {
	if params.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware)

	// auth
	r.POST("/signup", handlers.Auth.CreateUser)
	r.POST("/authenticate", handlers.Auth.Authenticate)

	// default public endpoints
	r.GET("/word/:word", handlers.Words.Word)
	r.GET("/words", handlers.Words.WordList)

	withAuth := r.Group("/").Use(middlewares.Auth)
	{
		// jar
		withAuth.GET("/jar", handlers.Jar.GetUserWords)
		withAuth.POST("/jar", handlers.Jar.AddWordToJar)

		// examine user
		withAuth.GET("/test", handlers.Exam.TestUser)
		withAuth.POST("/test", handlers.Exam.VerifyWord)

		// user
		withAuth.GET("/user", handlers.Auth.UserInfo)

	}

	admin := withAuth.Use(middlewares.IsAdmin)
	{
		admin.POST("/words", handlers.Words.AddWord)
	}

	return r
}
