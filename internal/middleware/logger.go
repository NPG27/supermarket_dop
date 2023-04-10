package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)

		status := c.Writer.Status()
		method := c.Request.Method

		if raw != "" {
			path = path + "?" + raw
		}

		logLine := fmt.Sprintf("%v | %3d | %13v | %15s |%-7s %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			status,
			latency,
			c.ClientIP(),
			method,
			path,
		)

		if _, err := logFile.WriteString(logLine); err != nil {
			log.Fatalf("Error writing log line to file: %v", err)
		}
	}
}
