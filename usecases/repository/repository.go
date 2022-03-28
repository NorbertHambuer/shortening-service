package repository

import (
	"github.com/norby7/shortening-service/usecases/repository/storage"
)

type Repository interface{
	storage.Storage
}
