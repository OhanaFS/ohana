package httpwc

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type HttpWriteCloser struct {
	wr *io.PipeWriter
	wg *sync.WaitGroup
}

func NewHttpWriteCloser(client *http.Client, method, url string) io.WriteCloser {
	// Set up pipes and a WaitGroup
	rd, wr := io.Pipe()
	wg := &sync.WaitGroup{}

	// Start the goroutine to upload the file
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Construct a new request
		req, err := http.NewRequest(method, url, rd)
		if err != nil {
			rd.CloseWithError(err)
			return
		}

		// Perform the request
		resp, err := client.Do(req)
		if err != nil {
			rd.CloseWithError(err)
			return
		}

		// Close with error if the response is not OK
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			rd.CloseWithError(fmt.Errorf("HTTP error: %d", resp.StatusCode))
			return
		}

		// Discard all output
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	// Return the writer
	return &HttpWriteCloser{wr, wg}
}

// Close closes the writer and waits for the upload to finish
func (wc *HttpWriteCloser) Close() error {
	// Close the pipe
	if err := wc.wr.Close(); err != nil {
		return err
	}

	// Wait for the upload to finish
	wc.wg.Wait()
	return nil
}

// Write writes the data to the writer
func (wc *HttpWriteCloser) Write(p []byte) (n int, err error) {
	return wc.wr.Write(p)
}
