package router

import (
	"github.com/eralves01/api_gateway/internal/middleware"
	"github.com/eralves01/api_gateway/internal/services"
	"github.com/eralves01/api_gateway/pkg/rate_limiter"
	"github.com/gorilla/mux"
)

func SetupRouter(limiter rate_limiter.Limiter) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.Logging)
	router.Use(middleware.RateLimitMiddleware(limiter))

	router.HandleFunc("/service/authenticate", services.ProxyRequest).Methods("GET", "POST", "PUT", "DELETE")

	return router
}
