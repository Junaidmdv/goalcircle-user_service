package redis

import (
	"github.com/junaidmdv/goalcirlcle/user_service/internal/config"
	redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(redisConfig *config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 10,
	})

	return rdb

}
