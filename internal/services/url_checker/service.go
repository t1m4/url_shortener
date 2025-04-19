package url_checker

import (
	"context"
	"time"
	"url_shortener/internal/services/api_client"
	"url_shortener/internal/utils"
)

type URLCheckerService struct {
	apiClient api_client.APIClient
}

func New() *URLCheckerService {
	return &URLCheckerService{
		apiClient: api_client.New(time.Minute),
	}
}

func (u *URLCheckerService) CheckURL(url string) ([]byte, error) {
	result, err := u.apiClient.Get(context.Background(), url)
	if err != nil {
		return nil, err
	}
	resultBytes, err := utils.PrettyString(result)
	return resultBytes, err
}
