package url_redirect

import "url_shortener/internal/repositories"

type URLRedirectService struct {
	shortenerRepository repositories.ShortenerRepository
}

func New(repositories *repositories.Repositories) *URLRedirectService {
	return &URLRedirectService{repositories.ShortenerRepository}
}

func (u *URLRedirectService) FindRedirectURL(url string) (string, error) {
	originalUrl, err := u.shortenerRepository.GetShortener(url)
	// Add logger logic
	return originalUrl, err
}
