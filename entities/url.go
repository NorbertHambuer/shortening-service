package entities

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"io"
	"net/url"
)

// Url defines the structure for the url object
// swagger: model
type Url struct {
	// the id for this url
	//
	// required: true
	Id int64 `json:"id"`
	// short url code
	//
	// min: 8
	// max: 8
	Code string `json:"code" validate:"required,min=8,max=8"`
	// original url
	//
	// min: 8
	Url string `json:"url" validate:"required,min=8"`
	// shortened url
	//
	// min: 16
	ShortUrl string `json:"shortUrl"`
	// shortened url domain
	//
	// min: 8
	Domain string `json:"domain" validate:"required,min=8"`
	// integer that represents the number of times the redirect endpoint was called for the code
	//
	// min: 0
	Counter int64 `json:"counter" validate:"gte=0"`
}

// Validate checks and validates each field of the Url object based on its definition
func (u *Url) Validate() error {
	validate := validator.New()

	_, err := url.Parse(u.Url)
	if err != nil {
		return err
	}

	_, err = url.Parse(u.ShortUrl)
	if err != nil {
		return err
	}

	return validate.Struct(u)
}

// ToJSON serializes the contents of the object to JSON
func (u *Url) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// FromJSON deserializes the JSON into the object
func (u *Url) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}
