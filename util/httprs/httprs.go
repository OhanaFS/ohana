package httprs

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/snabb/httpreaderat"
)

var (
	// ErrNoContentLength is returned when the server does not provide a
	// Content-Length header.
	ErrNoContentLength = errors.New("no Content-Length header")
	// ErrNoContentRange is returned when the server does not provide a
	// Content-Range header.
	ErrNoContentRange = errors.New("no Content-Range header")
	// ErrNoAcceptRanges is returned when the server does not provide an
	// Accept-Ranges header.
	ErrNoAcceptRanges = errors.New("no Accept-Ranges header")
)

// NewHttpRS creates a new HttpRS object.
func NewHttpRS(ctx context.Context, url string) (io.ReadSeeker, error) {
	// Initialize an HTTP client
	client := &http.Client{}

	// Send a HEAD request to get the size of the file.
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Make sure the server provies an Accept-Ranges header.
	if resp.Header.Get("Accept-Ranges") == "" {
		return nil, ErrNoAcceptRanges
	}

	// Get the size of the file.
	size := resp.Header.Get("Content-Length")
	if size == "" {
		return nil, ErrNoContentLength
	}
	sizeInt, err := strconv.ParseInt(size, 10, 64)

	// Initialize the HTTPReaderAt.
	getReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	hra, err := httpreaderat.New(client, getReq, nil)
	if err != nil {
		return nil, err
	}

	return io.NewSectionReader(hra, 0, sizeInt), nil
}
