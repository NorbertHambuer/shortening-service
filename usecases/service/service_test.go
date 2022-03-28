package service

import (
	"errors"
	"fmt"
	"github.com/norby7/shortening-service/entities"
	"testing"
)

var (
	addError     = fmt.Errorf("unable to add the url")
	deleteError  = fmt.Errorf("unable to delete the url")
	getError     = fmt.Errorf("unable to fetch the url")
	counterError = fmt.Errorf("unable to increment counter")
)

type RepositoryMock struct{}

func (r *RepositoryMock) Add(u *entities.Url) error {
	if u.Url == "http://www.invalidUrl.com" {
		return addError
	}

	return nil
}

func (r *RepositoryMock) Delete(id int64) error {
	if id == 0 {
		return deleteError
	}

	return nil
}

func (r *RepositoryMock) GetUrlByCode(code string) (string, error) {
	if code == "invalidCode" {
		return "", getError
	}

	if code != "84gfj4i9"{
		return "", nil
	}

	return "https://google.com", nil
}

func (r *RepositoryMock) GetById(id int64) (entities.Url, error) {
	if id == 0 {
		return entities.Url{}, getError
	}

	if id != 1{
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

func (r *RepositoryMock) GetByUrl(url string) (entities.Url, error) {
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

func (r *RepositoryMock) IncrementCounter(url string) error {
	if url == "" {
		return counterError
	}

	return nil
}

func TestDelete(t *testing.T) {
	r := &RepositoryMock{}
	s := NewService(r, 0, "http://localhost")

	testCases := []struct {
		name          string
		input         int64
		expectedError error
	}{
		{
			name:          "valid id, no error",
			input:         1,
			expectedError: nil,
		},
		{
			name:          "zero value id, error",
			input:         0,
			expectedError: deleteError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.Delete(tc.input)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got (%v)", tc.expectedError, err.Error())
			}
		})
	}
}

func TestGetUrlByCode(t *testing.T) {
	r := &RepositoryMock{}
	s := NewService(r, 0, "http://localhost")

	testCases := []struct {
		name          string
		input         string
		expectedError error
	}{
		{
			name:          "valid code, no error",
			input:         "4j83df92",
			expectedError: nil,
		},
		{
			name:          "invalid code, error",
			input:         "invalidCode",
			expectedError: getError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.GetUrlByCode(tc.input)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got (%v)", tc.expectedError, err.Error())
			}
		})
	}
}

func TestGetById(t *testing.T) {
	r := &RepositoryMock{}
	s := NewService(r, 0, "http://localhost")

	testCases := []struct {
		name          string
		input         int64
		expectedError error
	}{
		{
			name:          "valid id, no error",
			input:         1,
			expectedError: nil,
		},
		{
			name:          "zero value id, error",
			input:         0,
			expectedError: getError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.GetById(tc.input)

			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got (%v)", tc.expectedError, err.Error())
			}
		})
	}
}

func TestCreate(t *testing.T) {
	r := &RepositoryMock{}
	s := NewService(r, 0, "http://localhost")

	testCases := []struct {
		name    string
		input   *entities.Url
		isError bool
	}{
		{
			name:    "valid url, no code",
			input:   &entities.Url{Url: "www.validUrl.com"},
			isError: false,
		},
		{
			name:    "valid url, with code",
			input:   &entities.Url{Url: "http://www.validUrl.com", Code: "4j83df92"},
			isError: false,
		},
		{
			name:    "empty url",
			input:   &entities.Url{},
			isError: true,
		},
		{
			name:    "valid url, fetch error",
			input:   &entities.Url{Url: "http://www.invalidUrl.com"},
			isError: true,
		},
		{
			name:    "valid url, url already exists",
			input:   &entities.Url{Url: "www.existingUrl.com"},
			isError: false,
		},
		{
			name:    "no code,  code exists error",
			input:   &entities.Url{Url: "www.validUrl.com", Code: "invalidCode"},
			isError: true,
		},
		{
			name:    "with code, code exists error",
			input:   &entities.Url{Url: "http://www.validUrl.com", Code: "invalidCode"},
			isError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			err := s.Create(tc.input)

			if (err != nil) != tc.isError {
				t.Errorf("expected error (%v), got error (%v)", tc.isError, err)
			}
		})
	}
}
