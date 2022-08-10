package inc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/OhanaFS/ohana/dbfs"
	"github.com/OhanaFS/ohana/util"
	"github.com/OhanaFS/ohana/util/httprs"
	"github.com/OhanaFS/ohana/util/httpwc"
	"github.com/OhanaFS/stitch"
	"github.com/gorilla/mux"
)

// handleShardStream is the handler for /api/v1/node/shard/{shardId} defined in
// controller/inc/inc.go. It handles receiving of shards and storing them in the
// local filesystem, and also serving of shards from the local filesystem.
func (i *Inc) handleShardStream(w http.ResponseWriter, r *http.Request) {
	shardId := mux.Vars(r)["shardId"]
	localShardPath := path.Join(i.ShardsLocation, shardId)

	switch r.Method {
	case http.MethodHead:
	case http.MethodGet:
		// Try to open the file
		file, err := os.Open(localShardPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				util.HttpError(w, http.StatusNotFound, "Shard not found")
				return
			}
			util.HttpError(w, http.StatusInternalServerError, "Error opening file")
			return
		}
		defer file.Close()

		// Serve the file
		stat, err := file.Stat()
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, "Error getting file stats")
			return
		}
		http.ServeContent(w, r, shardId, stat.ModTime(), file)
		return
	case http.MethodPut:
		// Open the file for writing
		file, err := os.Create(localShardPath)
		if err != nil {
			util.HttpError(w, http.StatusInternalServerError, "Error opening file")
			return
		}
		defer file.Close()

		// Copy the file
		if _, err = io.Copy(file, r.Body); err != nil {
			util.HttpError(w, http.StatusInternalServerError, "Error copying file")
			return
		}

		// Flush file to disk
		file.Sync()

		// Finalize the shard headers
		enc := &stitch.Encoder{}
		if err := enc.FinalizeHeader(file); err != nil {
			util.HttpError(w, http.StatusInternalServerError, "Error finalizing header")
			return
		}

		// Return an OK
		util.HttpJson(w, http.StatusOK, "OK")
		return
	default:
		util.HttpError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
}

// getShardURL returns the URL of a shard on a server.
func (i *Inc) getShardURL(serverName, shardId string) (string, error) {
	// Get the address of the server
	host, err := dbfs.GetServerAddress(i.Db, serverName)
	if err != nil {
		return "", fmt.Errorf("failed to get server address: %v", err)
	}

	// Return the URL
	return fmt.Sprintf("https://%s/api/v1/node/shard/%s", host, shardId), nil
}

// NewShardWriter returns an io.WriteCloser to write a shard remotely to a server.
func (i *Inc) NewShardWriter(ctx context.Context, serverName, shardId string) (io.WriteCloser, error) {
	addr, err := i.getShardURL(serverName, shardId)
	if err != nil {
		return nil, err
	}

	wc := httpwc.NewHttpWriteCloser(ctx, i.HttpClient, http.MethodPut, addr)
	return wc, nil
}

// NewShardReader returns an io.ReadSeeker to read a shard remotely from a server.
func (i *Inc) NewShardReader(ctx context.Context, serverName, shardId string) (io.ReadSeeker, error) {
	addr, err := i.getShardURL(serverName, shardId)
	if err != nil {
		return nil, err
	}

	return httprs.NewHttpRS(ctx, i.HttpClient, addr)
}
