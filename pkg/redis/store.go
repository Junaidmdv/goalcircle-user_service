package redis

import (
	redisconfig "github.com/junaidmdv/goalcirlcle/authservice/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() *RedisStore {
	return &RedisStore{
		client: redisconfig.NewRedisClient(),
	}
}

func (rs *RedisStore) AddOtp() {

}  



