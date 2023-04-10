package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token header is missing"})
			return
		}
		expectedToken := os.Getenv("TOKEN")
		if token != expectedToken {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			return
		}
		ctx.Next()
	}
}
