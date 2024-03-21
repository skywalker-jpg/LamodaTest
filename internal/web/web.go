package web

import (
	"LamodaTest/internal/config"
	"LamodaTest/internal/storage"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"log/slog"
)

type Server struct {
	app     *echo.Echo
	URL     string
	logger  *slog.Logger
	Storage *storage.Storage
}

func New(srvCfg config.Server, logger *slog.Logger, storage *storage.Storage) (*Server, error) {
	e := echo.New()
	server := Server{
		app:     e,
		URL:     srvCfg.URL,
		logger:  logger,
		Storage: storage,
	}
	e.HideBanner = true
	e.Logger.SetOutput(io.Discard)

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	m := NewMiddleware(logger)
	m.Register(e)
	server.RegisterHandlers()

	return &server, nil
}

func (s *Server) Serve() error {
	s.logger.Info("HTTP server started", slog.String("url", s.URL))

	return fmt.Errorf("server error: %w", s.app.Start(s.URL))
}
