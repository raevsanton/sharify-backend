package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	serr := make(chan error, 1)

	go func() {
		fmt.Println("Server is listening on", s.httpServer.Addr)
		serr <- s.httpServer.ListenAndServe()
	}()

	var err error
	select {
	case err = <-serr:
	case <-ctx.Done():
	}

	sdctx, sdcancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer sdcancel()
	shutdownErr := s.httpServer.Shutdown(sdctx)

	if err != nil {
		return errors.Join(err, shutdownErr)
	}
	return shutdownErr
}
