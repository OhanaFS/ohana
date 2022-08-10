package httprs

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/snabb/httpreaderat"
)

// NewHttpRS creates a new HttpRS object.
func NewHttpRS(ctx context.Context, client *http.Client, url string) (io.ReadSeeker, error) {
	// Initialize the HTTPReaderAt.
	getReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	hra, err := httpreaderat.New(client, getReq, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating httpreaderat: %s", err)
	}

	return io.NewSectionReader(hra, 0, hra.Size()), nil
}
