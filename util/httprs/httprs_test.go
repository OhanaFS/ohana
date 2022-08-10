package httprs_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OhanaFS/ohana/util/httprs"
	"github.com/stretchr/testify/assert"
)

func TestHttpRS(t *testing.T) {
	assert := assert.New(t)

	sampleData := make([]byte, 1024)
	n, err := rand.Read(sampleData)
	assert.NoError(err)
	assert.Equal(1024, n)

	// Set up a mock HTTP server
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeContent(
					w, r, "test.txt", time.Now(),
					bytes.NewReader(sampleData),
				)
			},
		),
	)

	ctx := context.Background()
	hrs, err := httprs.NewHttpRS(ctx, &http.Client{}, server.URL+"/test.txt")
	assert.NoError(err)

	// Read
	buf := make([]byte, 32)
	n, err = hrs.Read(buf)
	assert.NoError(err)
	assert.Equal(32, n)
	assert.Equal(sampleData[:32], buf)

	// Seek to end of the file
	buf = make([]byte, 32)
	_, err = hrs.Seek(-32, io.SeekEnd)
	assert.NoError(err)
	n, err = hrs.Read(buf)
	assert.NoError(err)
	assert.Equal(32, n)
	assert.Equal(sampleData[len(sampleData)-32:], buf)
}
