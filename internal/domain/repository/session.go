package repository

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/redis/go-redis/v9"
)

type SessionStorage interface {
	SaveSession(context.Context, string, *entity.Session) error
}

type sessionStorage struct {
	redis *redis.Client
}

func NewSessionStorage(redis_client *redis.Client) SessionStorage {
	return &sessionStorage{
		redis: redis_client,
	}
}

func (rs *sessionStorage) SaveSession(ctx context.Context, sesionId string, sesion *entity.Session) error {
	res, err := rs.redis.HSet(ctx, sesionId, &sesion).Result()

	if err != nil {
		return err
	}

	t := reflect.TypeOf(sesion)

	if t.Len() != int(res) {
		return fmt.Errorf("some fields are missing")
	}

	return nil

}
 