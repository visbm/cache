package httpserver

import (
	"context"

	"cache/internal/config"
	"cache/pkg/logger"
	"net/http"
)

type Server struct {
	Logger     logger.Logger
	httpServer *http.Server
}

func NewServer(cnfg config.HttpServer, router http.Handler, logger logger.Logger) *Server {
	srv := &http.Server{
		Addr:           cnfg.Host + ":" + cnfg.Port,
		Handler:        router,
		MaxHeaderBytes: 1 << 20, //1 Mb
		ReadTimeout:    cnfg.ReadTimeout,
		WriteTimeout:   cnfg.WriteTimeout,
	}

	return &Server{
		Logger:     logger,
		httpServer: srv,
	}
}

func (s *Server) Start() error {

	s.Logger.Infof("Server starts at %s", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
