package initiator

import (
	"context"
	"time"

	"folderstructure/platform/logger"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func initRedis(redisUrl string, log logger.Logger) *redis.Client {
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatal(context.Background(), "unable to parse redis url", zap.Error(err))
	}

	opt.PoolSize = viper.GetInt("redis.pool_size")
	if opt.PoolSize == 0 {
		opt.PoolSize = 20
	}

	opt.IdleTimeout = viper.GetDuration("redis.idle_timeout")
	if opt.IdleTimeout == 0 {
		opt.IdleTimeout = 5 * time.Minute
	}

	opt.MinIdleConns = 5
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second

	rdb := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal(context.Background(), "Failed to connect to Redis", zap.Error(err))
	}

	log.Info(context.Background(), "Redis connected successfully", zap.String("addr", opt.Addr), zap.Int("pool_size", opt.PoolSize))

	return rdb
}
