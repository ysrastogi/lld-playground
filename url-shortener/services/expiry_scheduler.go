package services

import (
	"time"
	"url-shortener/src/repositories"
)

type ExpiryScheduler struct {
	UrlRepository repositories.URLRepository
}

func (e *ExpiryScheduler) schedule() {
	filters := map[string]interface{}{"expires_at": time.Now().Add(time.Hour * 24 * 90), "status": "ACTIVE"}
	ids, err := e.UrlRepository.ListURLs(filters)
	if err != nil {
		return
	}
	for _, id := range ids {
		e.UrlRepository.DeleteById(id.ID)
	}
}
