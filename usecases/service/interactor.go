package service

import "github.com/norby7/shortening-service/entities"

type Interactor interface {
	Create(*entities.Url) error
	Delete(int64) error
	GetUrlByCode(string) (string, error)
	GetById(int64) (entities.Url, error)
	IncrementCounter(string)
}
