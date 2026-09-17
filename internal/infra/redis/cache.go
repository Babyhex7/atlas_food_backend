package redis

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type CacheService struct {
	client  *Client
	sfGroup singleflight.Group
}

var (
	cacheServiceInstance *CacheService
	cacheOnce            sync.Once
)

func GetCacheService() *CacheService {
	cacheOnce.Do(func() {
		cacheServiceInstance = &CacheService{
			client: GetInstance(),
		}
	})
	return cacheServiceInstance
}

// GetOrFetch executes Cache-Aside pattern with SingleFlight to prevent cache stampede
func (cs *CacheService) GetOrFetch(ctx context.Context, key string, ttl time.Duration, target interface{}, fetchFunc func() (interface{}, error)) error {
	if cs == nil || !cs.client.IsAvailable() {
		// Fallback directly to DB fetch if Redis is unavailable
		res, err := fetchFunc()
		if err != nil {
			return err
		}
		return convertValue(res, target)
	}

	// 1. Try Read from Cache
	cachedVal, err := cs.client.Get(ctx, key)
	if err == nil && cachedVal != "" {
		if err := json.Unmarshal([]byte(cachedVal), target); err == nil {
			return nil
		}
	}

	// 2. Singleflight lock to deduplicate concurrent DB queries
	val, err, _ := cs.sfGroup.Do(key, func() (interface{}, error) {
		// Double check cache inside singleflight thread
		val, err := cs.client.Get(ctx, key)
		if err == nil && val != "" {
			var checkTarget interface{}
			if err := json.Unmarshal([]byte(val), &checkTarget); err == nil {
				return checkTarget, nil
			}
		}

		// Fetch from Source DB
		dbRes, fetchErr := fetchFunc()
		if fetchErr != nil {
			return nil, fetchErr
		}

		// Add random jitter (0 to 10% of TTL) to prevent simultaneous key expiration
		jitter := time.Duration(rand.Int63n(int64(ttl / 10)))
		dataBytes, marshalErr := json.Marshal(dbRes)
		if marshalErr == nil {
			if err := cs.client.Set(ctx, key, string(dataBytes), ttl+jitter); err != nil {
				log.Printf("⚠️ [CacheService] Failed to cache key %s: %v", key, err)
			}
		}

		return dbRes, nil
	})

	if err != nil {
		return err
	}

	return convertValue(val, target)
}

// InvalidateFoodCaches clears cached items related to a food or category
func (cs *CacheService) InvalidateFoodCaches(ctx context.Context, foodID string, categoryCode string) {
	if cs == nil || !cs.client.IsAvailable() {
		return
	}
	go func() {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if foodID != "" {
			cs.client.Del(ctxTimeout, fmt.Sprintf("atlas:prod:cache:foods:detail:%s", foodID))
		}
		if categoryCode != "" {
			cs.client.DelPattern(ctxTimeout, fmt.Sprintf("atlas:prod:cache:categories:%s:foods:*", categoryCode))
		}
		cs.client.Del(ctxTimeout, "atlas:prod:cache:categories:all")
		cs.client.DelPattern(ctxTimeout, "atlas:prod:cache:foods:search:*")
	}()
}

// InvalidateCategoryCaches clears all category caches
func (cs *CacheService) InvalidateCategoryCaches(ctx context.Context) {
	if cs == nil || !cs.client.IsAvailable() {
		return
	}
	go func() {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cs.client.Del(ctxTimeout, "atlas:prod:cache:categories:all")
		cs.client.DelPattern(ctxTimeout, "atlas:prod:cache:categories:*:foods:*")
	}()
}

// Helper to convert interface{} result to target struct via JSON
func convertValue(src interface{}, dest interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}

// GenerateSearchCacheKey builds a consistent MD5 hash cache key for search queries
func GenerateSearchCacheKey(query string, foodType string, limit int) string {
	raw := fmt.Sprintf("q=%s&type=%s&limit=%d", query, foodType, limit)
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("atlas:prod:cache:foods:search:%s", hex.EncodeToString(hash[:]))
}
