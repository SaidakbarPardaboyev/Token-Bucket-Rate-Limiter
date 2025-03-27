package main

import (
	"time"

	rl "github.com/SaidakbarPardaboyev/Token-Bucket-Rate-Limiter"
	"github.com/gin-gonic/gin"
)

func TestEndpoint(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "request recieved successfully",
	})
}

func main() {
	r := gin.Default()
	rateLimiter := rl.New()
	rateLimiter.SetConfig(rl.RateLimiterConfig{
		RATE_LIMIT:      100,
		REFILL_INTERVAL: time.Second,
	})
	rateLimiter.Run()

	r.Use(rateLimiter.RateLimitMiddleware())

	test := r.Group("/test")
	{
		test.GET("/test", TestEndpoint)
	}
}
