package httprs_test

import (
	"context"
	"embed"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OhanaFS/ohana/util/httprs"
	"github.com/stretchr/testify/assert"
)

//go:embed httprs_test.go
var testFS embed.FS

func TestHttpRS(t *testing.T) {
	assert := assert.New(t)

	// Set up a mock HTTP server
	fs := http.FileServer(http.FS(testFS))
	server := httptest.NewServer(fs)

	ctx := context.Background()
	hrs, err := httprs.NewHttpRS(ctx, server.URL+"/httprs_test.go")
	assert.NoError(err)

	// Read
	buf := make([]byte, 32)
	n, err := hrs.Read(buf)
	assert.NoError(err)
	assert.Equal(32, n)
	assert.True(string(buf[:8]) == "package ")

	// Seek to end of the file
	buf = make([]byte, 16)
	_, err = hrs.Seek(-16, io.SeekEnd)
	assert.NoError(err)
	n, err = hrs.Read(buf)
	assert.NoError(err)
	assert.Equal(16, n)
	assert.True(
		strings.HasSuffix(strings.TrimSpace(string(buf)), "END OF THE FILE"),
	)
}

// LEAVE THIS COMMENT AT THE END OF THE FILE
