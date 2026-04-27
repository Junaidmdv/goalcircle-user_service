package repository

import (
	"context"
	"time"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/internal/domain/entity"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type SessionStorage interface {
	SaveSession(context.Context, string, *entity.Session) error
	GetSession(context.Context, string) (*entity.Session, error)
	DeleteSession(context.Context, string) error
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

	pipe := rs.redis.Pipeline()
	pipe.HSet(ctx, key, *session)
	expiresAt, _ := time.Parse(time.RFC3339, session.ExpiresAt)
	pipe.ExpireAt(ctx, key, expiresAt)

	_, err := pipe.Exec(ctx)
	if err != nil {
		rs.logger.Error("failed store session data", "error", err, "data", session)
		return domain.NewInternalError("Something went wrong. Please try again later", err)
	}

	return nil
}

func (rs *sessionStorage) GetSession(ctx context.Context, key string) (*entity.Session, error) {
	var session entity.Session
	err := rs.redis.HGetAll(ctx, key).Scan(&session)
	if err != nil {
		rs.logger.Error("redis error", "error", err, "method", "GetSession")
		return nil, domain.NewInternalError("Something went wrong. Please try again later", err)
	}
	return &session, nil
}

func (rs *sessionStorage) DeleteSession(ctx context.Context, key string) error {
	deleted, err := rs.redis.Del(ctx, key).Result()
	if err != nil {
		rs.logger.Error("failed to delete session", "error", err, "key", key)
		return domain.NewInternalError("Something went wrong. Please try again later", err)
	}

	if deleted == 0 {
		return domain.NewNotFoundError("session not found")
	}

	return nil
}
 