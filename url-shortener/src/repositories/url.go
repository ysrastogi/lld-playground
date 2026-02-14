package repositories

import (
	"time"
	"url-shortener/src/models"

	"gorm.io/gorm"
)

type URLRepository interface {
	Save(url *models.URL, data map[string]interface{})
	FindByShortURL(shortCode string) (*models.URL, error)
	FindByLongURL(url string) (*models.URL, error)
	ListURLs(filters map[string]interface{}) ([]models.URL, error)
	Delete(shortCode string)
	DeleteById(id uint) (bool, error)
	Exists(shortCode string) (bool, error)
	Update(shortCode string, url *models.URL)
}

type URLRepositoryImpl struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return &URLRepositoryImpl{db: db}
}

func (r *URLRepositoryImpl) Save(url *models.URL, data map[string]interface{}) {
	shortCode := data["shortCode"].(string)
	longURL := data["url"].(string)
	isStale := false
	status := models.ACTIVE
	expiresAt := time.Now().Add(time.Hour * 24 * 90)

	r.db.Create(&models.URL{
		ShortCode: shortCode,
		Url:       longURL,
		IsStale:   isStale,
		Status:    status,
		ExpiresAt: expiresAt,
	})
}

func (r *URLRepositoryImpl) FindByShortURL(shortCode string) (*models.URL, error) {
	var url models.URL
	if err := r.db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepositoryImpl) FindByLongURL(longURL string) (*models.URL, error) {
	var result models.URL
	if err := r.db.Where("url = ?", longURL).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *URLRepositoryImpl) Delete(shortCode string) {
	var url models.URL
	if err := r.db.Where("short_code = ? AND is_stale = ?", shortCode, false).First(&url).Error; err != nil {
		return
	}
	r.db.Model(&url).Update("is_stale", true)
}

func (r *URLRepositoryImpl) DeleteById(id uint) (bool, error) {
	if err := r.db.Model(&models.URL{}).Where("id = ?", id).Update("is_stale", true).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *URLRepositoryImpl) Exists(shortCode string) (bool, error) {
	var url models.URL
	if err := r.db.Where("short_code = ? AND is_stale = ?", shortCode, false).First(&url).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *URLRepositoryImpl) Update(shortCode string, url *models.URL) {
	var existing models.URL
	if err := r.db.Where("short_code = ? AND is_stale = ?", shortCode, false).First(&existing).Error; err != nil {
		return
	}
	r.db.Model(&existing).Update("url", url.Url)
}

func (r *URLRepositoryImpl) ListURLs(filters map[string]interface{}) ([]models.URL, error) {
	var urls []models.URL
	if err := r.db.Where(filters).Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}
