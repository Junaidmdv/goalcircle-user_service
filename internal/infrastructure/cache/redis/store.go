package redis

import (
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() *RedisStore {
	return &RedisStore{
		client: NewRedisClient(),
	}
}
func (rs *RedisStore) AddOtp() {

}  



