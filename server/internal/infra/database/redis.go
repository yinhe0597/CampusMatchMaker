package database

import (
	"context"
	"fmt"

	"campus_collab/internal/infra/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// InitRedis 初始化 Redis 连接
func InitRedis(cfg config.RedisConfig, log *zap.Logger) (*redis.Client, error) {
	if !cfg.Enabled {
		log.Info("Redis 缓存已禁用")
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Redis 连接失败: %w", err)
	}

	log.Info("Redis 连接成功",
		zap.String("addr", cfg.Addr()),
		zap.Int("db", cfg.DB),
	)
	return client, nil
}
