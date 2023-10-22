package web

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// LoggingMiddleware is a middleware function for create `zerolog` log entries for all Gin Gonic handler method executions.
func LoggingMiddleware(g *gin.Context) {
	start := time.Now()
	g.Next()
	elapsed := time.Since(start)

	log.Info().
		Str("method", g.Request.Method).
		Str("URI", g.Request.RequestURI).
		Int("status", g.Writer.Status()).
		Dur("latency", elapsed).
		Str("client", g.ClientIP()).
		Msg("GIN REQUEST")
	g.Next()
}
