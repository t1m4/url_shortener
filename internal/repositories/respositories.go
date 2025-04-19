package repositories

import "gorm.io/gorm"

type Repositories struct {
	ShortenerRepository ShortenerRepository
}

func New(db *gorm.DB) *Repositories {
	return &Repositories{
		ShortenerRepository: NewShorternerRepository(db),
	}
}
