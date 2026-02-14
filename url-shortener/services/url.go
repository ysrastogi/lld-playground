package services

import (
	"errors"
	"url-shortener/src/models"
	"url-shortener/src/repositories"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}

type RateLimiter interface {
	IsAllowed(key string, limitPolicy string) (bool, error)
}
type AnalyticsPublisher interface {
	Publish(metrics *models.Metrics) error
}

type URLService interface {
	CreateShortURL(longURL string, userId string) (string, error)
	Redirect(shortCode string) (string, error)
	Deactivate(shortCode string) error
}

type URLServiceImpl struct {
	urlRepository      repositories.URLRepository
	cache              Cache
	shortCodeStrategy  ShortCodeStrategy
	rateLimiter        RateLimiter
	analyticsPublisher AnalyticsPublisher
}

func NewURLService(
	urlRepository repositories.URLRepository,
	cache Cache,
	shortCodeStrategy ShortCodeStrategy,
	rateLimiter RateLimiter,
	analyticsPublisher AnalyticsPublisher,
) URLService {
	return &URLServiceImpl{
		urlRepository:      urlRepository,
		cache:              cache,
		shortCodeStrategy:  shortCodeStrategy,
		rateLimiter:        rateLimiter,
		analyticsPublisher: analyticsPublisher,
	}
}

func (s *URLServiceImpl) CreateShortURL(longURL string, userId string) (string, error) {
	allowed, err := s.rateLimiter.IsAllowed(userId, "create")
	if err != nil {
		return "", err
	}
	if !allowed {
		return "", errors.New("rate limit exceeded")
	}

	existing, err := s.urlRepository.FindByLongURL(longURL)
	if err == nil && existing != nil {
		return existing.ShortCode, nil
	}

	var shortCode string
	for {
		shortCode = s.shortCodeStrategy.generate(longURL)
		exists, err := s.urlRepository.Exists(shortCode)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
	}

	url := &models.URL{}
	data := map[string]interface{}{
		"shortCode": shortCode,
		"url":       longURL,
	}
	s.urlRepository.Save(url, data)

	_ = s.cache.Set(shortCode, longURL)

	return shortCode, nil
}

func (s *URLServiceImpl) Redirect(shortCode string) (string, error) {
	longURL, err := s.cache.Get(shortCode)
	if err == nil && longURL != "" {
		_ = s.analyticsPublisher.Publish(&models.Metrics{ShortCode: shortCode})
		return longURL, nil
	}

	url, err := s.urlRepository.FindByShortURL(shortCode)
	if err != nil {
		return "", errors.New("short URL not found")
	}

	if url.Status != models.ACTIVE {
		return "", errors.New("short URL is not active")
	}

	_ = s.cache.Set(shortCode, url.Url)

	_ = s.analyticsPublisher.Publish(&models.Metrics{ShortCode: shortCode})

	return url.Url, nil
}

func (s *URLServiceImpl) Deactivate(shortCode string) error {
	s.urlRepository.Delete(shortCode)
	_ = s.cache.Delete(shortCode)
	return nil
}
