package redis

import (
	"fmt"

	"github.com/Junaidmdv/goalcircle-user_service/internal/config"
	redis "github.com/redis/go-redis/v9"
) 


type Redis struct{
	Client *redis.Client
}

func NewRedisClient(redisConfig *config.RedisConfig) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       0,
		PoolSize: 10,
	})

	return &Redis{
       Client: rdb,
	}

}
