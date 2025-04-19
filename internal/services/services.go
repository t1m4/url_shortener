package services

import (
	"time"
	"url_shortener/configs"
	"url_shortener/internal/repositories"
	"url_shortener/internal/services/api_client"
	"url_shortener/internal/services/url_checker"
	"url_shortener/internal/services/url_redirect"
	"url_shortener/internal/services/url_shortener"
)

type Service struct {
	apiClient           api_client.APIClient
	URLCheckerService   *url_checker.URLCheckerService
	URLShortenerService *url_shortener.URLShortenerService
	URLRedirectService  *url_redirect.URLRedirectService
}

func New(config *configs.Config, repositories *repositories.Repositories) *Service {
	url_checker_service := url_checker.New()
	apiClient := api_client.New(time.Minute)
	return &Service{
		apiClient:           apiClient,
		URLCheckerService:   url_checker_service,
		URLShortenerService: url_shortener.New(config, repositories, apiClient),
		URLRedirectService:  url_redirect.New(repositories),
	}
}
