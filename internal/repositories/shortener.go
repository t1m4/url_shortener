package repositories

import (
	"errors"
	"fmt"
	"url_shortener/internal/custom_errors"
	"url_shortener/internal/db"
	"url_shortener/internal/logger"
	"url_shortener/internal/schemas"

	"gorm.io/gorm"
)

type ShortenerRepository interface {
	InsertShortener(*schemas.URLInput, string) (*db.Shortener, error)
	GetShortener(string) (string, error)
}
type shortenerRepository struct {
	l  logger.Logger
	db *gorm.DB
}

func NewShorternerRepository(l logger.Logger, db *gorm.DB) ShortenerRepository {
	return &shortenerRepository{l, db}
}

func (s *shortenerRepository) InsertShortener(urlInput *schemas.URLInput, newLink string) (*db.Shortener, error) {
	tx := s.db.Begin()
	var shortener db.Shortener
	result := tx.Select("id", "new_link").Where("link = ?", urlInput.Url).First(&shortener)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		shortener = db.Shortener{Link: urlInput.Url, NewLink: newLink}
		result := tx.Create(&shortener)
		if result.Error != nil {
			tx.Commit()
			s.l.Error(result.Error)
			return nil, fmt.Errorf(custom_errors.DbError)
		}
	}
	tx.Commit()
	return &shortener, nil
}

func (s *shortenerRepository) GetShortener(newLink string) (string, error) {
	tx := s.db.Begin()
	var shortener db.Shortener
	result := tx.Select("id", "link").Where("new_link = ?", newLink).First(&shortener)
	if result.Error != nil {
		tx.Commit()
		s.l.Error(result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf(custom_errors.RowDoesNotExistError)
		}
		return "", fmt.Errorf(custom_errors.DbError)
	}
	tx.Model(&shortener).Update("count", gorm.Expr("count + 1")).Where("new_link = ?", newLink)
	tx.Commit()
	return shortener.Link, nil
}

type FakeShortenerRepository struct {
	InsertResult *db.Shortener
	InsertErr    error
	GetResult    string
	GetErr       error
}

func (f *FakeShortenerRepository) InsertShortener(*schemas.URLInput, string) (*db.Shortener, error) {
	return f.InsertResult, f.InsertErr
}
func (f *FakeShortenerRepository) GetShortener(string) (string, error) {
	return f.GetResult, f.GetErr
}
