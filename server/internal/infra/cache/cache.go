package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache Redis 缓存抽象层
type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewCache 创建缓存实例
// client 为 nil 时缓存功能静默禁用（不报错）
func NewCache(client *redis.Client, ttl time.Duration) *Cache {
	return &Cache{client: client, ttl: ttl}
}

// IsEnabled 检查缓存是否可用
func (c *Cache) IsEnabled() bool {
	return c != nil && c.client != nil
}

// Get 从缓存获取数据，反序列化到 dest
// 缓存未命中返回 ErrCacheMiss
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	if c == nil || c.client == nil {
		return ErrCacheMiss
	}

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return ErrCacheMiss
	}
	if err != nil {
		return fmt.Errorf("缓存读取失败: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		_ = c.client.Del(ctx, key)
		return fmt.Errorf("缓存反序列化失败: %w", err)
	}
	return nil
}

// Set 将数据序列化后写入缓存
func (c *Cache) Set(ctx context.Context, key string, value interface{}) error {
	if c == nil || c.client == nil {
		return nil
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("缓存序列化失败: %w", err)
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}

// Delete 删除指定缓存键
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

// DeletePattern 按模式删除缓存键（用于批量失效）
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	if c == nil || c.client == nil {
		return nil
	}

	iter := c.client.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}
