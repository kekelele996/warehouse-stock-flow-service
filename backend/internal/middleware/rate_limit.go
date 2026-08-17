package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/wmsflow/wmsflow/internal/constants"
	"github.com/wmsflow/wmsflow/internal/util"
)

// RateLimit 限流中间件：优先使用 Redis 固定窗口计数，Redis 不可用时降级到进程内计数。
func RateLimit(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	fallback := newMemoryLimiter(limit, window)
	return func(c *gin.Context) {
		key := "rl:" + c.ClientIP()
		if redisClient != nil {
			allowed, err := redisAllow(c.Request.Context(), redisClient, key, limit, window)
			if err == nil {
				if !allowed {
					util.Fail(c, http.StatusTooManyRequests, constants.CodeTooManyRequests, "请求过于频繁，请稍后再试")
					c.Abort()
					return
				}
				c.Next()
				return
			}
		}
		if !fallback.allow(key) {
			util.Fail(c, http.StatusTooManyRequests, constants.CodeTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}

func redisAllow(ctx context.Context, client *redis.Client, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now().UnixMilli()
	start := now - window.Milliseconds()
	pipe := client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", fmtInt(start))
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmtInt(now)})
	pipe.Expire(ctx, key, window)
	countCmd := pipe.ZCard(ctx, key)
	if _, err := pipe.Exec(ctx); err != nil {
		return false, err
	}
	return countCmd.Val() <= int64(limit), nil
}

func fmtInt(v int64) string {
	return time.UnixMilli(v).Format("20060102150405.000")
}

type memoryLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	hits   map[string][]int64
}

func newMemoryLimiter(limit int, window time.Duration) *memoryLimiter {
	return &memoryLimiter{limit: limit, window: window, hits: make(map[string][]int64)}
}

func (m *memoryLimiter) allow(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().Unix()
	cutoff := now - int64(m.window.Seconds())
	list := m.hits[key]
	kept := list[:0]
	for _, t := range list {
		if t >= cutoff {
			kept = append(kept, t)
		}
	}
	if len(kept) >= m.limit {
		m.hits[key] = kept
		return false
	}
	m.hits[key] = append(kept, now)
	return true
}
