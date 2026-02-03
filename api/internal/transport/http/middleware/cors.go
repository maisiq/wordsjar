package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var AllowedHosts = []string{
	"http://frontend.local.test:3000",
	"localhost",
	"null",
}

func CORSMiddleware(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")

	ctx.SetSameSite(http.SameSiteLaxMode)

	origin := ctx.GetHeader("Origin")

	for _, o := range AllowedHosts {
		if origin == o {
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PUTCH, DELETE, OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Credentials", "true")
			break
		}
	}

	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
