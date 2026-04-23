package repository

import (
	"context"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type SessionStorage interface {
	SaveSession(context.Context, string, *entity.Session) error
}

type sessionStorage struct {
	redis  *redis.Client
	logger logger.Logger
}

func NewSessionStorage(redis_client *redis.Client) SessionStorage {
	return &sessionStorage{
		redis: redis_client,
	}
}

func (rs *sessionStorage) SaveSession(ctx context.Context, key string, session *entity.Session) error {

	_, err := rs.redis.HSet(ctx, key, session).Result()
	if err != nil {
		rs.logger.Error("failed store session data", "error", err, "data", session)
		return domain.NewInternalError("Something went wrong.Please try again later", err)
	}

	return nil

}
