package service

import (
	"fmt"
	"github.com/norby7/shortening-service/entities"
	"github.com/norby7/shortening-service/usecases/repository"
	"log"
	"math/rand"
	"strings"
)

type Service struct {
	Repo        repository.Repository
	CounterJobs chan string
	Domain string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// NewService returns a new Service object address
func NewService(r repository.Repository, workers int, domain string) *Service {
	counterJobs := make(chan string, 100)
	for i := 0; i < workers; i++ {
		go counterWorker(r, counterJobs)
	}

	return &Service{Repo: r, CounterJobs: counterJobs, Domain: domain}
}

// counterWorker fetches short url codes from a channel and calls the repository IncrementCounter function with the code
func counterWorker(repo repository.Repository, jobs <-chan string) {
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}

			err := repo.IncrementCounter(job)
			if err != nil {
				log.Printf("unable to increment code (%s) counter: %s\n", job, err.Error())
			}
		}
	}
}

// Create validates the Url object, generates a new code if none is given and inserts it into the repository
func (s *Service) Create(u *entities.Url) error {
	// check if the url has no scheme
	if !strings.HasPrefix(u.Url, "http://") && !strings.HasPrefix(u.Url, "https://") {
		u.Url = "http://" + u.Url
	}

	// check if the url exists, return the shortUrl if it does
	dbUrl, err := s.Repo.GetByUrl(u.Url)
	if err != nil{
		return fmt.Errorf("unable to check if url already exist in the database: %s", err.Error())
	}

	// if the url is found, return it
	if dbUrl.Id != 0{
		*u = dbUrl
		return nil
	}

	// if no code was sent by the user, generate a new unique code
	if u.Code == "" {
		code, err := s.generateNewUniqueCode()
		if err != nil {
			return err
		}

		u.Code = code
	} else {
		// check if the code already exists
		exists, err := s.codeExists(u.Code)
		if err != nil {
			return fmt.Errorf("%s: %s", ErrCheckCode.Error(), err.Error())
		}

		if exists {
			return ErrCodeAlreadyExists
		}
	}

	u.ShortUrl = s.Domain + "/" + u.Code
	u.Domain = s.Domain

	// validate the Url object
	if err := u.Validate(); err != nil {
		return err
	}

	return s.Repo.Add(u)
}

// Delete removes a Url from the repository
func (s *Service) Delete(id int64) error {
	return s.Repo.Delete(id)
}

// GetUrlByCode fetches a Url from the repository by its code
func (s *Service) GetUrlByCode(code string) (string, error) {
	return s.Repo.GetUrlByCode(code)
}

// GetById fetches a Url from the repository by its id
func (s *Service) GetById(id int64) (entities.Url, error) {
	return s.Repo.GetById(id)
}

// IncrementCounter adds a new code into the CounterJobs channel
func (s *Service) IncrementCounter(code string) {
	s.CounterJobs <- code
}

// randCode returns a random string with n length
func randCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// codeExists checks if the code is already stored into the database
func (s *Service) codeExists(code string) (bool, error) {
	// check if code already exists
	url, err := s.Repo.GetUrlByCode(code)
	if err != nil {
		return false, fmt.Errorf("%s: %s", ErrCheckCode.Error(), err.Error())
	}

	return url != "", nil
}

// generateNewUniqueCode creates a new code
// if the generated code is not unique, it will regenerate it
func (s *Service) generateNewUniqueCode() (string, error) {
	code := randCode(8)
	// while the code already exists
	for {
		// check if code already exists
		exists, err := s.codeExists(code)
		if err != nil {
			return "", fmt.Errorf("%s: %s", ErrCheckCode.Error(), err.Error())
		}

		if !exists {
			break
		}

		code = randCode(8)
	}

	return code, nil
}