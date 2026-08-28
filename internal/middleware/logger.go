package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)
		event := log.Info()

		switch {
		case status >= 500:
			event = log.Error()
		case status >= 400:
			event = log.Warn()
		}

		event.
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP()).
			Msg("request processed")
	}
}
