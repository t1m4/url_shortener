package url_shortener

import (
	"context"
	"fmt"
	"strings"
	"url_shortener/configs"
	"url_shortener/internal/custom_errors"
	"url_shortener/internal/repositories"
	"url_shortener/internal/schemas"
	"url_shortener/internal/services/api_client"
	"url_shortener/internal/utils"
)

type URLShortenerService struct {
	config              *configs.Config
	shortenerRepository repositories.ShortenerRepository
	getNewLink          func(int) string
	apiClient           api_client.APIClient
}

func New(config *configs.Config, repositories *repositories.Repositories, apiClient api_client.APIClient) *URLShortenerService {
	return &URLShortenerService{config, repositories.ShortenerRepository, utils.GetNewLink, apiClient}
}

func (u *URLShortenerService) ShortURL(urlInput *schemas.URLInput) (string, error) {
	if urlInput.Url == "" {
		return "", fmt.Errorf(custom_errors.UrlRequiredError)
	}
	urlInput.Url = strings.Trim(urlInput.Url, " ")
	newLink := u.getNewLink(configs.UrlLenght)
	_, err := u.apiClient.Get(context.Background(), urlInput.Url)
	if err != nil {
		return "", fmt.Errorf(custom_errors.MakingRequestError)
	}

	shorneter, err := u.shortenerRepository.InsertShortener(urlInput, newLink)
	if err != nil {
		return "", err
	}
	newUrl := strings.Join([]string{u.config.APP.DOMAIN, "/api/", shorneter.NewLink}, "")
	return newUrl, nil
}
