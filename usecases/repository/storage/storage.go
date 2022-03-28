package storage

import "github.com/norby7/shortening-service/entities"

type Storage interface{
	Add(*entities.Url) error
	Delete(int64) error
	GetUrlByCode(string) (string, error)
	GetById(int64) (entities.Url, error)
	GetByUrl(string) (entities.Url, error)
	IncrementCounter(string) error
}

