package http

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/jasonhancock/go-logger"
)

// NewHTTPServer starts up an HTTP server. The server will run until the context
// is cancelled.
func NewHTTPServer(ctx context.Context, l *logger.L, wg *sync.WaitGroup, hler http.Handler, addr string) error {
	server := http.Server{
		Addr:    addr,
		Handler: hler,
	}

	wg.Add(1)
	go func() {
		l.Info("starting http server", "addr", addr)
		// TODO: add TLS support
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Err("starting http server error", "error", err.Error(), "addr", addr)
		}
	}()

	go func() {
		defer wg.Done()
		<-ctx.Done()
		l.Info("stopping http server", "addr", addr)

		// shut down gracefully, but wait no longer than 10 seconds before halting
		sdCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(sdCtx); err != nil {
			l.Err(
				"stopping http server",
				"error", err.Error,
				"addr", addr,
			)
		}
		l.Info("stopped http server", "addr", addr)
	}()

	return nil
}
