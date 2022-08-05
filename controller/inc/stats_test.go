package inc_test

import (
	"fmt"
	"github.com/OhanaFS/ohana/controller/inc"
	"github.com/OhanaFS/ohana/util/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerSpecifics(t *testing.T) {

	t.Run("Getting df output test", func(t *testing.T) {

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

		report, err := i.GetLocalServerStatusReport()
		assert.NoError(t, err)
		assert.NotNil(t, report)
		fmt.Println("Report:", report)

	})

}
