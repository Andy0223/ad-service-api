package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Before request
		log.Printf("Started %s %s", c.Request.Method, c.Request.RequestURI)

		c.Next() // 处理请求

		// After request
		latency := time.Since(t)
		log.Printf("Completed in %v", latency)
	}
}
