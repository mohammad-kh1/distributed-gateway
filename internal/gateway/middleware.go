package gateway

import (
	"net/http"
	"time"

	"github.com/mohammad-kh1/distributed-gateway/internal/logger"
	"go.uber.org/zap"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)
		duration := time.Since(start)
		logger.Info("Inbound Request", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("remote_addr", r.RemoteAddr), zap.Duration("lantecy", duration))

	})

}
