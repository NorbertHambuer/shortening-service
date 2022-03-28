// Package http, classification of url shortening service API
//
// Documentation for url shortening service API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package http

import (
	"encoding/json"
	"fmt"
	"github.com/norby7/shortening-service/entities"
	"github.com/norby7/shortening-service/usecases/service"
	"log"
	"net/http"
	"path"
	"strconv"
)

// Data structure representing a single url
// swagger:response urlResponse
type urlResponse struct {
	// Location header that contains the short url
	Location string
	// A single url object
	// in: body
	Body entities.Url
}

// Generic error message response
// swagger:response errorResponse
type errorResponse struct {
	Message string `json:"message"`
}

// Code already exists in the database error message response
// swagger:response codeExistsErrorResponse
type codeExistsErrorResponse struct {
	Message string `json:"message"`
}

// swagger:response noContent
type noContent struct{}

// Redirections counter response
// swagger:response counterResponse
type counterResponse struct {
	Counter int64 `json:"counter"`
}

// swagger:model
type addParam struct {
	// short url code
	//
	// required: false
	// min: 8
	// max: 8
	Code string `json:"code" validate:"required,min=8,max=8"`
	// original url
	//
	// required: true
	// min: 8
	Url string `json:"url" validate:"required,min=8"`
}

// swagger:parameters Add
type urlParam struct {
	// Url object used for Add<br>
	// Note: the "code" field is optional, the service will generate a random code if none given
	// in: body
	// required: true
	Body addParam
}

// swagger:parameters Delete Get GetCounter
type Id struct {
	// Url object Id
	// in: path
	// required: true
	Id int64
}

// swagger:parameters Redirect
type Code struct {
	// Url object Code
	// in: path
	// required: true
	Code string
}

type Controller struct {
	Service service.Interactor
	Logger  *log.Logger
}

func NewController(s service.Interactor, l *log.Logger) *Controller {
	return &Controller{Service: s, Logger: l}
}

// swagger:route POST /api api Add
// Creates a new url in the database and then returns it in the response
// responses:
// 201: urlResponse
// 409: codeExistsErrorResponse
// 422: errorResponse
// 500: errorResponse

// Add creates a new url in the database and returns it
func (c *Controller) Add(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	c.Logger.Println("Handle Add url")

	var u entities.Url
	err := u.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": "unable to parse url object %s"}`, err.Error()), http.StatusUnprocessableEntity)
		return
	}

	if err = c.Service.Create(&u); err != nil {
		code := http.StatusInternalServerError
		msg := err
		if err == service.ErrCodeAlreadyExists {
			code = http.StatusConflict
			msg = service.ErrCodeAlreadyExists
		}

		http.Error(rw, fmt.Sprintf(`{"message": "unable to add url %s"}`, msg.Error()), code)
		return
	}

	rw.Header().Set("Location", u.ShortUrl)
	rw.WriteHeader(http.StatusCreated)
	if err = u.ToJSON(rw); err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": "unable to encode url response object %s"}`, err.Error()), http.StatusUnprocessableEntity)
		return
	}
}

// swagger:route DELETE /api/{Id} api Delete
// Deletes a url
// responses:
// 200: noContent
// 400: errorResponse
// 500: errorResponse

// Delete removes a url from the database
func (c *Controller) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	c.Logger.Println("Handle Delete url")

	id, err := strconv.Atoi(path.Base(r.URL.String()))
	if err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": "invalid url id value: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if err := c.Service.Delete(int64(id)); err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": "unable to delete url: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/{Id} api Get
// Returns a url based on the given ID or 404 if no short url exists with the given code
// responses:
// 200: urlResponse
// 404: noContent
// 422: errorResponse
// 500: errorResponse

// Get fetches a url from the database
func (c *Controller) Get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	c.Logger.Println("Handle get url")

	id, err := strconv.Atoi(path.Base(r.URL.String()))
	if err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": "invalid url id value: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	url, err := c.Service.GetById(int64(id))
	if err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": "unable to fetch url: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	if url.Id == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if err = url.ToJSON(rw); err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": "unable to encode url response object %s"}`, err.Error()), http.StatusUnprocessableEntity)
		return
	}
}

// swagger:route GET /{Code} root Redirect
// Redirects to a long url or returns 404 if no short url exists in the database with the given code
// responses:
// 302: noContent
// 404: noContent
// 500: errorResponse

// RedirectShortUrl redirects the request to a long url if the given code exists in the database
func (c *Controller) RedirectShortUrl(rw http.ResponseWriter, r *http.Request) {
	c.Logger.Println("Handle url redirect")

	code := path.Base(r.URL.String())
	url, err := c.Service.GetUrlByCode(code)
	if err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": "unable to fetch url: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	if url == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	c.Service.IncrementCounter(code)

	http.Redirect(rw, r, url, http.StatusFound)
}

// swagger:route GET /counter/{Id} counter GetCounter
// Returns the redirections counter for a given url object Id
// responses:
// 200: counterResponse
// 404: noContent
// 500: errorResponse

// GetCounter returns the redirections counter for a given url object Id
func (c *Controller) GetCounter(rw http.ResponseWriter, r *http.Request) {
	c.Logger.Println("Handle get counter")

	id, err := strconv.Atoi(path.Base(r.URL.String()))
	if err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": "invalid url id value: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	url, err := c.Service.GetById(int64(id))
	if err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": ""unable to fetch url: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	if url.Id == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if err = json.NewEncoder(rw).Encode(counterResponse{Counter: url.Counter}); err != nil {
		http.Error(rw, fmt.Sprintf(`{"message": ""unable to encode url response object %s"}`, err.Error()), http.StatusUnprocessableEntity)
		return
	}

}
