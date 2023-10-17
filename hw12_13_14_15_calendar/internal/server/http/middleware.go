package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type ResponseHook struct {
	http.ResponseWriter
	statusCode int
}

func (r *ResponseHook) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.statusCode = statusCode
}

func loggingMiddleware(next http.Handler, log Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		rh := &ResponseHook{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rh, r)
		latency := time.Since(now).Milliseconds()

		log.Info(fmt.Sprintf(
			"%s [%s] %s %s %d %d %s",
			r.Host,
			now.String(),
			r.Method,
			r.URL.Path,
			rh.statusCode,
			latency,
			r.UserAgent(),
		))
	})
}
