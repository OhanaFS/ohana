package service

import (
	"context"
	"github.com/OhanaFS/ohana/config"
	"github.com/OhanaFS/ohana/service/kv"
	"github.com/OhanaFS/ohana/util"
	"strconv"
	"time"
)

// Session is the interface for the kv service. It is used to create,
// get, and invalidate sessions. It can use different backends to store
// sessions.
//
// Sessions are key-value pairs, where the key is the kv id and the value
// is the user id. The caller is responsible for fetching the user data from
// the database.
type Session interface {
	Get(ctx context.Context, sessionId string) (uint, error)
	Create(ctx context.Context, userId uint, ttl time.Duration) (string, error)
	Invalidate(ctx context.Context, key string) error
}

type session struct {
	store kv.KV
}

func NewSession(cfg *config.Config) (Session, error) {

	return &session{
		store: kv.NewRedis(cfg),
	}, nil
}

func (s *session) Get(ctx context.Context, sessionId string) (uint, error) {
	struid, err := s.store.Get(ctx, sessionId)
	if err != nil {
		return 0, err
	}

	uid, err := strconv.ParseUint(struid, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(uid), nil
}

func (s *session) Create(ctx context.Context, userId uint, ttl time.Duration) (string, error) {
	// Generate a new kv id
	id, err := util.RandomHex(32)
	if err != nil {
		return "", err
	}

	return id, s.store.Set(ctx, id, strconv.FormatUint(uint64(userId), 10), ttl)
}

func (s *session) Invalidate(ctx context.Context, key string) error {
	return s.store.Delete(ctx, key)
}
