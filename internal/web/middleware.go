package web

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"log/slog"
	"strconv"
	"time"
)

type Middleware struct {
	logger *slog.Logger
}

func NewMiddleware(logger *slog.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) Register(router *echo.Echo) {
	router.Use(m.AccessLog())
}

func (m *Middleware) AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()

			requestID := uuid.New().String()
			c.Set("requestID", requestID)

			m.logger.Info("Request started",
				slog.String("RequestID", requestID),
				slog.String("IP", c.RealIP()),
				slog.String("URL", c.Request().URL.Path),
				slog.String("Method", c.Request().Method),
			)

			err := next(c)

			responseTime := time.Since(startTime)
			if err != nil {
				m.logger.Error("Request Failed",
					slog.String("RequestID", requestID),
					slog.String("Time spent", strconv.FormatInt(int64(responseTime), 10)),
					slog.String("Error", err.Error()),
				)
			} else {
				m.logger.Info("Request done",
					slog.String("RequestID", requestID),
					slog.String("Time spent", strconv.FormatInt(int64(responseTime), 10)),
				)
			}

			return err
		}
	}
}
