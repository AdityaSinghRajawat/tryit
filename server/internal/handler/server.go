package handler

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/config"
	"github.com/AdityaSinghRajawat/tryit/server/internal/utils/logger"
)

// Server wraps the loopback bind + timeouts + graceful drain (§9.1).
type Server struct {
	cfg config.Config
	log *logger.Logger
	hs  *http.Server
}

func NewServer(cfg config.Config, log *logger.Logger, h http.Handler) *Server {
	return &Server{
		cfg: cfg,
		log: log,
		hs: &http.Server{
			Addr:              cfg.ListenAddr(),
			Handler:           h,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      cfg.ExecTimeout + 30*time.Second, // accommodate slow targets
			IdleTimeout:       60 * time.Second,
		},
	}
}

// ListenAndServe binds to 127.0.0.1 only (per G2).
func (s *Server) ListenAndServe() error {
	ln, err := net.Listen("tcp", s.hs.Addr)
	if err != nil {
		return err
	}
	s.log.Info("listening", "addr", s.hs.Addr)
	if err := s.hs.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error { return s.hs.Shutdown(ctx) }
