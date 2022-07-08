package service

import (
	"context"
	"time"

	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/service/kv"
	"github.com/OhanaFS/ohana/util"
	"go.uber.org/zap"
)

// Session is the interface for the kv service. It is used to create,
// get, and invalidate sessions. It can use different backends to store
// sessions.
//
// Sessions are key-value pairs, where the key is the kv id and the value
// is the user id. The caller is responsible for fetching the user data from
// the database.
type Session interface {
	Get(ctx context.Context, sessionId string) (string, error)
	Create(ctx context.Context, userId string, ttl time.Duration) (string, error)
	Invalidate(ctx context.Context, key string) error
}

type session struct {
	store kv.KV
}

func NewSession(cfg *config.Config, logger *zap.Logger) (Session, error) {
	if cfg.Redis.Address == "" {
		logger.Warn("No redis address configured, using in-memory kv for session storage")
		return &session{
			store: kv.NewMemoryKV(),
		}, nil
	}

	return &session{
		store: kv.NewRedis(cfg),
	}, nil
}

func (s *session) Get(ctx context.Context, sessionId string) (string, error) {
	uid, err := s.store.Get(ctx, string(sessionId))
	if err != nil {
		return "", err
	}

	return uid, nil
}

func (s *session) Create(ctx context.Context, userId string, ttl time.Duration) (string, error) {
	// Generate a new kv id
	id, err := util.RandomHex(32)
	if err != nil {
		return "", err
	}

	return id, s.store.Set(ctx, id, string(userId), ttl)
}

func (s *session) Invalidate(ctx context.Context, key string) error {
	return s.store.Delete(ctx, string(key))
}
