package httpwc_test

import (
	"crypto/rand"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OhanaFS/ohana/util/httpwc"
	"github.com/stretchr/testify/assert"
)

func TestHttpWC(t *testing.T) {
	assert := assert.New(t)

	// Generate test data
	data := make([]byte, 1024)
	n, err := rand.Read(data)
	assert.NoError(err)
	assert.Equal(len(data), n)

	// Create a mock HTTP server to receive the uploads
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Received request: %s %s", r.Method, r.URL.Path)
		inData, err := io.ReadAll(r.Body)
		t.Logf("Received data: %d bytes", len(inData))
		assert.NoError(err)
		assert.Equal(data, inData)
	}))
	defer server.Close()

	// Create a client
	t.Logf("Creating client")
	writer := httpwc.NewHttpWriteCloser(&http.Client{}, "POST", server.URL)

	// Write the data
	t.Logf("Writing data")
	n, err = writer.Write(data)
	assert.NoError(err)
	assert.Equal(len(data), n)

	// Close the writer
	t.Logf("Closing writer")
	err = writer.Close()
	assert.NoError(err)
}
