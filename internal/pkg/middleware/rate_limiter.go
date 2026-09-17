package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	redisInfra "atlas_food/internal/infra/redis"
)

// Lua script for atomic sliding window rate limiting using Redis Sorted Set (ZSet)
const slidingWindowLuaScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local windowMs = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local clearBefore = now - windowMs

redis.call('ZREMRANGEBYSCORE', key, 0, clearBefore)
local currentRequests = redis.call('ZCARD', key)

if currentRequests < limit then
    redis.call('ZADD', key, now, now)
    redis.call('EXPIRE', key, math.ceil(windowMs / 1000) + 1)
    return {1, currentRequests + 1}
else
    return {0, currentRequests}
end
`

// RateLimiter returns a Gin middleware that enforces sliding window rate limits per client IP
func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	redisClient := redisInfra.GetInstance()

	return func(c *gin.Context) {
		if redisClient == nil || !redisClient.IsAvailable() {
			// Graceful degradation: if Redis is offline, allow request
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		key := fmt.Sprintf("atlas:prod:ratelimit:%s:%s", clientIP, path)
		nowMs := time.Now().UnixMilli()
		windowMs := window.Milliseconds()

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		res, err := redisClient.RawClient().Eval(ctx, slidingWindowLuaScript, []string{key}, nowMs, windowMs, limit).Result()
		if err != nil {
			// If Redis query fails, allow request to proceed (fail open)
			c.Next()
			return
		}

		resSlice, ok := res.([]interface{})
		if !ok || len(resSlice) < 2 {
			c.Next()
			return
		}

		allowed := resSlice[0].(int64) == 1
		count := resSlice[1].(int64)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, int64(limit)-count)))

		if !allowed {
			c.Header("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Rate limit exceeded. Please wait before retrying.",
			})
			return
		}

		c.Next()
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
