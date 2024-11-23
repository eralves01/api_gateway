package middleware

import (
	"net/http"

	"github.com/eralves01/api_gateway/pkg/rate_limiter"
)

func RateLimitMiddleware(limiter rate_limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Identificador do cliente (pode ser o IP ou um token exclusivo)
			clientID := r.RemoteAddr

			// Verifica se a requisição é permitida
			if !limiter.Allow(clientID) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Requisição permitida, continue para o próximo handler
			next.ServeHTTP(w, r)
		})
	}
}
