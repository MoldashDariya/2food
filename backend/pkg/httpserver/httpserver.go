package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	closeTimeout = 5 * time.Second
)

type Config struct {
	Port         string        `default:"8080" envconfig:"HTTP_PORT"`
	ReadTimeout  time.Duration `default:"5s" envconfig:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `default:"35s" envconfig:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `default:"120s" envconfig:"HTTP_IDLE_TIMEOUT"`
}

type Server struct {
	httpServer *http.Server
}

func New(handler http.Handler, c Config) *Server {
	s := &Server{
		httpServer: &http.Server{
			Addr:         ":" + c.Port,
			Handler:      handler,
			ReadTimeout:  c.ReadTimeout,
			WriteTimeout: c.WriteTimeout,
			IdleTimeout:  c.IdleTimeout,
		},
	}
	return s
}

func (s *Server) Start() {
	go s.httpServer.ListenAndServe()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), closeTimeout)
	defer cancel()
	
	return s.httpServer.Shutdown(ctx)
}
