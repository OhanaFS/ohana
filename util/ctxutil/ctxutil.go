package ctxutil

import (
	"context"

	"github.com/OhanaFS/ohana/dbfs"
	"gorm.io/gorm"
)

const (
	ContextKeyTransaction = "transaction"
	ContextKeyUser        = "user"
)

// GetTransaction returns a transaction from the context if it exists. If a
// transaction does not exist, it returns the supplied database connection with
// its context changed to the supplied context.
func GetTransaction(ctx context.Context, db *gorm.DB) *gorm.DB {
	if v := ctx.Value(ContextKeyTransaction); v != nil {
		if val, ok := v.(*gorm.DB); ok {
			return val
		}
	}
	return db.WithContext(ctx)
}

// WithTransaction returns a new context with the given transaction.
func WithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, ContextKeyTransaction, tx)
}

// GetUser returns a user from the context if it exists. If a user does not
// exist, it returns an error.
func GetUser(ctx context.Context) (*dbfs.User, error) {
	if v := ctx.Value(ContextKeyUser); v != nil {
		if val, ok := v.(*dbfs.User); ok {
			return val, nil
		}
	}
	return &dbfs.User{}, dbfs.ErrCredentials
}

// WithUser returns a new context with the given user.
func WithUser(ctx context.Context, user *dbfs.User) context.Context {
	return context.WithValue(ctx, ContextKeyUser, user)
}
