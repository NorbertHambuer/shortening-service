package repository

import (
	"fmt"
	"github.com/norby7/shortening-service/entities"
	"log"
	"os"
	"testing"
)

var (
	addError     = fmt.Errorf("unable to add the url")
	deleteError  = fmt.Errorf("unable to delete the url")
	getError     = fmt.Errorf("unable to fetch the url")
	counterError = fmt.Errorf("unable to increment counter")
	getUrlError  = fmt.Errorf("unable to get url from cache")
	setUrlError  = fmt.Errorf("unable to get save url into cache")
)

type StorageMock struct{}
type CacheMock struct{}

func (r *StorageMock) Add(u *entities.Url) error {
	if u.Url == "http://www.invalidUrl.com" {
		return addError
	}

	return nil
}

func (r *StorageMock) Delete(id int64) error {
	if id == 0 {
		return deleteError
	}

	return nil
}

func (r *StorageMock) GetUrlByCode(code string) (string, error) {
	if code == "invalidCode" {
		return "", getError
	}

	if code != "84gfj4i9" && code != "invalidSetCode" {
		return "", nil
	}

	return "https://google.com", nil
}

func (r *StorageMock) GetById(id int64) (entities.Url, error) {
	if id == 0 {
		return entities.Url{}, getError
	}

	if id != 1 {
		return entities.Url{}, nil
	}

	return entities.Url{
		Id:       1,
		Code:     "84gfj4i9",
		Url:      "https://google.com",
		ShortUrl: "http://localhost/84gfj4i9",
		Domain:   "http://localhost",
		Counter:  1,
	}, nil
}

func (r *StorageMock) GetByUrl(url string) (entities.Url, error) {
	if url == "http://www.invalidUrl.com" {
		return entities.Url{}, getError
	}

	if url != "http://www.existingUrl.com" {
		return entities.Url{}, nil
	}

	return entities.Url{
		Id:       1,
		Code:     "84gfj4i9",
		Url:      "https://google.com",
		ShortUrl: "http://localhost/84gfj4i9",
		Domain:   "http://localhost",
		Counter:  1,
	}, nil
}

func (r *StorageMock) IncrementCounter(url string) error {
	if url == "" {
		return counterError
	}

	return nil
}

func (c *CacheMock) SetShortUrl(code, url string) error {
	if code == "invalidSetCode" {
		return setUrlError
	}

	return nil
}

func (c *CacheMock) GetShortUrl(code string) (string, error) {
	if code == "invalidCode" {
		return "", getUrlError
	}

	if code == "cacheUrl" {
		return "cacheUrl", nil
	}

	return "", nil
}

func TestGetUrlByCode(t *testing.T) {
	l := log.New(os.Stdout, "urls-api-test", log.LstdFlags)
	st := &StorageMock{}
	ch := &CacheMock{}
	repo := NewUrlRepository(st, ch, l)

	testCases := []struct {
		name    string
		input   string
		isError bool
	}{
		{
			name:    "empty code",
			input:   "",
			isError: false,
		},
		{
			name:    "invalid code",
			input:   "invalidCode",
			isError: true,
		},
		{
			name:    "save to cache error",
			input:   "invalidCode",
			isError: true,
		},
		{
			name:    "valid url from cache",
			input:   "cacheUrl",
			isError: false,
		},
		{
			name:    "valid url",
			input:   "84gfj4i9",
			isError: false,
		},
		{
			name:    "set url into cache error",
			input:   "invalidSetCode",
			isError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.GetUrlByCode(tc.input)

			if (err != nil) != tc.isError {
				t.Errorf("expected error (%v), got error (%v)", tc.isError, err)
			}
		})
	}
}
