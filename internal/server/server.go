package server

import (
	"fmt"
	"forum/internal/config"
	"forum/internal/service"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.ServerCfg, wcfg *config.WebCfg, srvc *service.Service) error {
	handler, err := newMainHandler(srvc, wcfg)
	if err != nil {
		return err
	}
	s.httpServer = &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler.InitRoutes(wcfg),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Printf("Server runs on http://localhost%s\n", s.httpServer.Addr)
	err = s.httpServer.ListenAndServe()
	return fmt.Errorf("Run: %w", err)
}
