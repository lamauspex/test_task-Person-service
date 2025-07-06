package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggerMiddleware возвращает Echo-middleware для введения журнала
func LoggerMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()

			// Следующая ступень обработчиков
			err := next(c)

			endTime := time.Now()
			responseTime := endTime.Sub(startTime)

			logger.Info(
				"HTTP Request",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Path()),
				zap.Int("status_code", c.Response().Status),
				zap.Duration("response_time", responseTime),
				zap.String("remote_ip", c.RealIP()),
			)

			return err
		}
	}
}
