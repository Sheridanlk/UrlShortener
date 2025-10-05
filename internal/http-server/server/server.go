package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	log *slog.Logger
	srv *http.Server
}

func New(log *slog.Logger, router http.Handler, addres string, readTimeout time.Duration, writeTimeout time.Duration, idleTimeout time.Duration) *Server {

	srv := &http.Server{
		Addr:         addres,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return &Server{
		log: log,
		srv: srv,
	}
}

func (s *Server) Run() error {
	const op = "http-server.server.Run"

	log := s.log.With(slog.String("op", op))
	log.Info("http-server is runnung", slog.String("addr", s.srv.Addr))

	return s.srv.ListenAndServe()
}

func (s *Server) Stop() {
	const op = "http-server.server.Stop"

	s.log.With(slog.String("op", op)).Info("Stopping http server")
	s.srv.Shutdown(context.Background())
}
