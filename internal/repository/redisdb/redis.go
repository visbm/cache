package redisdb

import (
	"cache/internal/config"	
	"cache/pkg/logger"
	"github.com/go-redis/redis/v7"
	"time"
)

type RedisDatabase struct {
	Client  *redis.Client
	timeout time.Duration
	logger  logger.Logger
}

func NewRedisDatabase(cfg config.Redis, logger logger.Logger) (*RedisDatabase, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisDatabase{
		Client:  rdb,
		logger:  logger,
		timeout: cfg.Timeout,
	}, nil
}

func (r *RedisDatabase) Close() error {
	return r.Client.Close()
}
