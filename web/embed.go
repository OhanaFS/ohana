package web

import (
	"embed"
	"io/fs"
)

// webApp is the embedded filesystem containing the built web application
//go:embed dist/*
var webApp embed.FS

// GetWebApp returns the embedded web application filesystem
func GetWebApp() (fs.FS, error) {
	return fs.Sub(webApp, "dist")
}
