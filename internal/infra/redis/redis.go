package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb       *redis.Client
	available bool
	mu        sync.RWMutex
}

var (
	instance *Client
	once     sync.Once
)

// GetInstance returns the singleton Redis client instance
func GetInstance() *Client {
	once.Do(func() {
		instance = NewRedisClient()
	})
	return instance
}

func NewRedisClient() *Client {
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")
	dbStr := getEnv("REDIS_DB", "0")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		db = 0
	}

	opts := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     password,
		DB:           db,
		PoolSize:     50,
		MinIdleConns: 10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   3,
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := &Client{
		rdb:       rdb,
		available: false,
	}

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️ [Redis] Connection failed to %s:%s - %v (Graceful fallback to DB enabled)", host, port, err)
	} else {
		client.available = true
		log.Printf("✅ [Redis] Connected successfully to %s:%s (DB: %d)", host, port, db)
	}

	return client
}

func (c *Client) RawClient() *redis.Client {
	return c.rdb
}

func (c *Client) IsAvailable() bool {
	if c == nil || c.rdb == nil {
		return false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.available
}

func (c *Client) SetAvailable(status bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.available = status
}

// Get helper with fail-safe fallback
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	if !c.IsAvailable() {
		return "", redis.Nil
	}
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		log.Printf("⚠️ [Redis] Error getting key %s: %v", key, err)
	}
	return val, err
}

// Set helper with TTL
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !c.IsAvailable() {
		return nil
	}
	err := c.rdb.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Printf("⚠️ [Redis] Error setting key %s: %v", key, err)
	}
	return err
}

// Delete helper
func (c *Client) Del(ctx context.Context, keys ...string) error {
	if !c.IsAvailable() || len(keys) == 0 {
		return nil
	}
	return c.rdb.Del(ctx, keys...).Err()
}

// DelPattern helper for prefix invalidation
func (c *Client) DelPattern(ctx context.Context, pattern string) error {
	if !c.IsAvailable() {
		return nil
	}
	iter := c.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) > 0 {
		return c.rdb.Del(ctx, keys...).Err()
	}
	return nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
