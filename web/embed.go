package web

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"strings"

	"github.com/OhanaFS/ohana/util"
)

// webApp is the embedded filesystem containing the built web application
//go:embed dist/*
var webApp embed.FS

// GetWebApp returns the embedded web application filesystem
func GetWebApp() (fs.FS, error) {
	return fs.Sub(webApp, "dist")
}

// GetHandler returns the handler which serves the web app
func GetHandler() (http.HandlerFunc, error) {
	webApp, err := GetWebApp()
	if err != nil {
		return nil, err
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to open the file
		path := r.URL.Path[1:]
		file, err := webApp.Open(path)
		if err != nil {
			// If error, return index.html
			path = "index.html"
			if file, err = webApp.Open(path); err != nil {
				util.HttpError(w, http.StatusInternalServerError,
					fmt.Sprintf("Failed to open index.html: %s", err.Error()))
				return
			}
		}

		// Guess mime type
		parts := strings.Split(path, ".")
		ext := "." + parts[len(parts)-1]
		mimeType := mime.TypeByExtension(ext)

		// Send headers
		w.Header().Add("Content-Type", mimeType)

		// Write the file
		io.Copy(w, file)
	})

	return handler, nil
}
