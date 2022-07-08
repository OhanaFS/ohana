package middleware

import (
	"net/http"

	"github.com/OhanaFS/ohana/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type middleware = func(next http.Handler) http.Handler
type Middlewares struct {
	Logging  middleware
	UserAuth middleware
}

func Provide(session service.Session, db *gorm.DB, logger *zap.Logger) *Middlewares {
	return &Middlewares{
		Logging:  NewLoggingMW(logger),
		UserAuth: NewUserAuthMW(session, db),
	}
}
