package middleware

import (
	"net/http"

	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/service"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"gorm.io/gorm"
)

const SessionCookieName = "ohana-session"

// NewUserAuthMW creates a middleware that checks if the user is logged in based
// on the session cookie. If the user is logged in, it updates the request
// context with the user and passes the request to the next handler.
//
// If the session is not found, the user is not found in the database, or the
// user is not allowed to log in (i.e. email is not verified, or the user has
// been suspended), the request is aborted with a StatusUnauthorized response.
//
// This middleware only performs authentication. The service is responsible for
// authorization.
func NewUserAuthMW(session service.Session, db *gorm.DB) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get the session sessionId
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				util.HttpError(w, http.StatusUnauthorized, "Session cookie not found")
				return
			}

			// Fetch the user ID from the session store
			uid, err := session.Get(ctx, cookie.Value)
			if err != nil {
				util.HttpError(w, http.StatusUnauthorized,
					"Session not found or has expired")
				return
			}

			// Fetch the user from the database
			tx := ctxutil.GetTransaction(ctx, db)
			user, err := dbfs.GetUserById(tx, uid)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Set the user in the request context
			r = r.WithContext(ctxutil.WithUser(ctx, user))

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
