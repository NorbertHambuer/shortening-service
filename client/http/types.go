package http

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"io"
	"net/url"
)

type Url struct {
	Id int64 `json:"id"`
	Code string `json:"code" validate:"required,min=8,max=8"`
	Url string `json:"url" validate:"required,min=8"`
	ShortUrl string `json:"shortUrl"`
	Domain string `json:"domain" validate:"required,min=8"`
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

type CreateRequest struct{
	Url string `json:"url" validate:"required,min=8"`
	Code string `json:"code"`
}

// ToJSON serializes the contents of the object to JSON
func (c *CreateRequest) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

// Validate checks and validates each field of the CreateRequest object based on its definition
func (c *CreateRequest) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}

type ErrorResponse struct{
	Message string `json:"message"`
}

type CounterResponse struct{
	Value int64 `json:"counter"`
}

// FromJSON deserializes the JSON into the object
func (cr *CounterResponse) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(cr)
}
