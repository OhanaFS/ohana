package inc_test

import (
	"context"
	"fmt"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/controller/middleware"
	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util/ctxutil"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServerSpecifics(t *testing.T) {

	t.Run("Getting Server Status Report", func(t *testing.T) {

		// init db

		db := testutil.NewMockDB(t)

		i := inc.Inc{
			ServerName:     "blahServerName",
			HostName:       "blahHostName",
			Port:           "5555",
			ShardsLocation: ".",
			BindIp:         "192.168.2.53",
			Db:             db,
		}

		// Testing with a valid name

		report, err := i.GetServerStatusReport("blahServerName")
		assert.NoError(t, err)
		assert.NotNil(t, report)
		fmt.Println("Report:", report)

		// Testing with an invalid name

		report, err = i.GetServerStatusReport("idkServer")
		assert.Error(t, err)
		assert.Nil(t, report)
		fmt.Println("Report:", report)

		// Testing with http request

		session := testutil.NewMockSession(t)
		sessionId, err := session.Create(nil, "superuser", time.Hour)

		user, err := dbfs.GetUser(db, "superuser")
		assert.NoError(t, err)

		req := httptest.NewRequest("GET", "/server/api/v1/stats/", nil).WithContext(
			ctxutil.WithUser(context.Background(), user))
		req.AddCookie(&http.Cookie{Name: middleware.SessionCookieName, Value: sessionId})
		w := httptest.NewRecorder()
		i.ReturnServerDetails(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		fmt.Println("Response:", w.Body.String())

		// Testing with an invalid name

	})

}
