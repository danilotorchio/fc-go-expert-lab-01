package transport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Server struct {
	http *http.Server
}

func (s *Server) Start(ctx context.Context, shutdownTimeout time.Duration) error {
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("HTTP server listening", "addr", s.http.Addr, "pid", strconv.Itoa(os.Getpid()))
		if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		slog.Info("termination signal received; shutting down server", "signal", ctx.Err())
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		if closeErr := s.http.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}

		return fmt.Errorf("shutting down http server: %w", err)
	}

	slog.Info("HTTP server shutdown gracefully")
	return nil
}

type NewServerOpts struct {
	Port    string
	Handler http.Handler
}

func NewServer(opts *NewServerOpts) *Server {
	return &Server{
		http: &http.Server{
			Addr:    net.JoinHostPort("", opts.Port),
			Handler: opts.Handler,

			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}
}
