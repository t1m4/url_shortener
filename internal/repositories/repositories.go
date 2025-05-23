package repositories

import (
	"url_shortener/internal/logger"

	"gorm.io/gorm"
)

type Repositories struct {
	ShortenerRepository ShortenerRepository
}

func New(l logger.Logger, db *gorm.DB) *Repositories {
	return &Repositories{
		ShortenerRepository: NewShorternerRepository(l, db),
	}
}
