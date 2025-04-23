package services

import (
	"time"
	"url_shortener/configs"
	"url_shortener/internal/logger"
	"url_shortener/internal/repositories"
	"url_shortener/internal/services/api_client"
	"url_shortener/internal/services/rate_limiter"
	"url_shortener/internal/services/url_checker"
	"url_shortener/internal/services/url_redirect"
	"url_shortener/internal/services/url_shortener"
)

type Service struct {
	apiClient           api_client.APIClient
	URLCheckerService   *url_checker.URLCheckerService
	URLShortenerService *url_shortener.URLShortenerService
	URLRedirectService  *url_redirect.URLRedirectService
	RateLimiterService  rate_limiter.RateLimiterService
}

func New(config *configs.Config, l logger.Logger, repositories *repositories.Repositories) *Service {
	url_checker_service := url_checker.New()
	apiClient := api_client.New(time.Minute)
	rateLimiterService := rate_limiter.NewRateLimiterService(config, l)
	return &Service{
		apiClient:           apiClient,
		URLCheckerService:   url_checker_service,
		URLShortenerService: url_shortener.New(config, repositories, apiClient),
		URLRedirectService:  url_redirect.New(repositories),
		RateLimiterService:  rateLimiterService,
	}
}

func (s *Service) Start() {
	s.RateLimiterService.Start()
}
func (s *Service) Stop() {
	s.RateLimiterService.Stop()
}
