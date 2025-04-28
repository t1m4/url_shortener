package middleware

import (
	"net/http"
	"url_shortener/configs"
	"url_shortener/internal/custom_errors"
	"url_shortener/internal/logger"
	"url_shortener/internal/services"
	"url_shortener/internal/services/rate_limiter"
)

type Middleware struct {
	config      *configs.Config
	l           logger.Logger
	rateLimiter rate_limiter.RateLimiterService
}

func New(config *configs.Config, l logger.Logger, services *services.Service) *Middleware {
	return &Middleware{config, l, services.RateLimiterService}
}
func (m *Middleware) CheckRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("user_id")
		err := m.rateLimiter.Check(userId)
		if err != nil {
			custom_errors.Write400(m.l, err.Error(), w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
