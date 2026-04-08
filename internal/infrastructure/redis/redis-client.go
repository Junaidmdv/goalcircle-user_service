package redis

import (
	"fmt"

	"github.com/junaidmdv/goalcirlcle/user_service/internal/config"
	redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(redisConfig *config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       0,
		PoolSize: 10,
	})

	return rdb

}
