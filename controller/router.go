package controller

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/OhanaFS/ohana/config"
)

// NewRouter creates a new gorilla/mux router.
func NewRouter() *mux.Router {
	return mux.NewRouter()
}

// StartServer creates a new HTTP server with the given router. It uses fx
// lifecycle hooks to start and stop the server.
func NewServer(
	lc fx.Lifecycle,
	router *mux.Router,
	logger *zap.Logger,
	config *config.Config,
) error {
	logger.Info(
		"Starting HTTP server",
		zap.String("address", config.HTTP.Bind),
	)

	// Set up the SPA router
	if config.SPA.UseDevelopmentServer {
		rpURL, err := url.Parse(config.SPA.DevelopmentServerURL)
		if err != nil {
			return err
		}
		rp := httputil.NewSingleHostReverseProxy(rpURL)
		router.PathPrefix("/").Handler(rp)
	} else {
		spa := &spaHandler{
			staticPath: config.SPA.StaticPath,
			indexPath:  config.SPA.IndexPath,
		}
		router.PathPrefix("/").Handler(spa)
	}

	// Wrap the router in a handler that logs requests
	handler := NewLoggingMiddleware(logger)(router)

	// Set up the server
	srv := &http.Server{
		Addr:    config.HTTP.Bind,
		Handler: handler,

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	// Create the lifecycle hook that starts and stops the server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, time.Second*15)
			defer cancel()

			return srv.Shutdown(ctx)
		},
	})

	return nil
}
