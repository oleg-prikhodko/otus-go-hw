package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(logger Logger, next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		latency := time.Since(start).Milliseconds()

		logger.Info(
			fmt.Sprintf(
				"[%s] %s %s %s %s %d %dms %s",
				start.Format(time.RFC3339), // Datetime
				r.RemoteAddr,               // Client IP
				r.Method,                   // Method
				r.URL.Path,                 // Path
				r.Proto,                    // HTTP Version
				rw.statusCode,              // Response Code
				latency,                    // Latency in ms
				r.UserAgent(),              // User Agent
			),
		)
	})
}
