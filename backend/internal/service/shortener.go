package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/DragonZero000/url-shortener/internal/model"
)

type urlRepository interface {
	Save(code, originalURL string) error
	GetByCode(code string) (*model.URL, error)
	GetByOriginalURL(originalURL string) (*model.URL, error)
}

type Service struct {
	db urlRepository
}

func NewService(repo urlRepository) *Service {
	return &Service{db: repo}
}

var ErrInvalidURL = errors.New("invalid URL")

func validateURL(rawURL string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return ErrInvalidURL
	}
	return nil
}

func (s *Service) Shorten(originalURL string) (string, error) {
	err := validateURL(originalURL)
	if err != nil {
		return "", err
	}
	existing, err := s.db.GetByOriginalURL(originalURL)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return existing.Code, nil
	}
	maxAttempts := 6
	for i := 0; i < maxAttempts; i++ {
		saltedURL := originalURL
		if i != 0 {
			saltedURL += strconv.Itoa(i)
		}
		code := sha256.Sum256([]byte(saltedURL))
		newCode := hex.EncodeToString(code[:])[:10]
		existing, err = s.db.GetByCode(newCode)
		if err != nil {
			return "", err
		}
		if existing == nil {
			err = s.db.Save(newCode, originalURL)
			if err != nil {
				return "", err
			}
			return newCode, nil
		}
	}
	return "", errors.New("hex collision")
}

var ErrNotFound = errors.New("not found")

func (s *Service) GetByCode(code string) (string, error) {
	existing, err := s.db.GetByCode(code)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return existing.OriginalURL, nil
	}
	return "", ErrNotFound
}
