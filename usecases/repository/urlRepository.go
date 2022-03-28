package repository

import (
	"github.com/go-redis/redis"
	"github.com/norby7/shortening-service/entities"
	"github.com/norby7/shortening-service/usecases/repository/cache"
	"github.com/norby7/shortening-service/usecases/repository/storage"
	"log"
)

type UrlRepository struct {
	storage storage.Storage
	cache   cache.Cache
	Logger  *log.Logger
}

// NewUrlRepository returns a new UrlRepository object address
func NewUrlRepository(s storage.Storage, c cache.Cache, l *log.Logger) *UrlRepository {
	return &UrlRepository{
		storage: s,
		cache:   c,
		Logger:  l,
	}
}

// Add calls the storage Add function to insert a new Url into the database
func (r *UrlRepository) Add(u *entities.Url) error {
	return r.storage.Add(u)
}

// Delete calls the storage Delete function to remove a Url from the database
func (r *UrlRepository) Delete(id int64) error {
	return r.storage.Delete(id)
}

// GetByUrlCode returns a long url either from the cache if it exists or from the storage if it doesn't
// It adds the code to the cache if it doesn't already exists
func (r *UrlRepository) GetUrlByCode(code string) (string, error) {
	// search code in cache
	u, err := r.cache.GetShortUrl(code)
	if err != nil && err != redis.Nil{
		r.Logger.Println("unable to get short url from cache: " + err.Error())
	}

	// if the code doesn't exist in cache
	if u == "" {
		// get url from storage
		url, err := r.storage.GetUrlByCode(code)
		if err != nil {
			return "", err
		}

		// if url exists, add it to the cache
		if url != "" {
			err = r.cache.SetShortUrl(code, url)
			if err != nil {
				r.Logger.Println("unable to add short url to cache:" + err.Error())
			}
		}

		return url, nil
	}

	return u, nil
}

// GetById calls the storage GetById function to fetch a Url from the database by its Id
func (r *UrlRepository) GetById(id int64) (entities.Url, error) {
	return r.storage.GetById(id)
}

// GetByUrl calls the storage GetByUrl function to fetch a Url from the database by its Url
func (r *UrlRepository) GetByUrl(url string) (entities.Url, error) {
	return r.storage.GetByUrl(url)
}

// IncrementCounter calls the storage IncrementCounter
func (r *UrlRepository) IncrementCounter(code string) error {
	return r.storage.IncrementCounter(code)
}
