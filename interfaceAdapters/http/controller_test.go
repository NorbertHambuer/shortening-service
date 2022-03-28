package http

import (
	"fmt"
	"github.com/norby7/shortening-service/entities"
	"github.com/norby7/shortening-service/usecases/service"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	createError  = fmt.Errorf("unable to create the url")
	deleteError  = fmt.Errorf("unable to delete the url")
	getError     = fmt.Errorf("unable to fetch the url")
)

type ServiceMock struct{}

func (s *ServiceMock) Create(u *entities.Url) error {
	if u.Url == "http://www.invalidUrl.com" || u.Url == "" {
		return createError
	}

	if u.Code == "d4jn8dsf" {
		return service.ErrCodeAlreadyExists
	}

	return nil
}

func (s *ServiceMock) Delete(id int64) error {
	if id == 0 {
		return deleteError
	}

	return nil
}

func (s *ServiceMock) GetUrlByCode(code string) (string, error) {
	if code == "invalidCode" {
		return "", getError
	}

	if code != "84gfj4i9"{
		return "", nil
	}

	return "https://google.com", nil
}

func (s *ServiceMock) GetById(id int64) (entities.Url, error) {
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

func (s *ServiceMock) IncrementCounter(string) {

}

func TestAdd(t *testing.T) {
	s := ServiceMock{}
	l := log.New(os.Stdout, "urls-api", log.LstdFlags)
	c := NewController(&s, l)

	testCases := []struct {
		name       string
		input      *strings.Reader
		statusCode int
	}{{
		name:       "invalid json object",
		input:      strings.NewReader(`"url":"Where does the sun set?}`),
		statusCode: http.StatusUnprocessableEntity,
	}, {
		name:       "add error",
		input:      strings.NewReader(`{"url":"http://www.invalidUrl.com"}`),
		statusCode: http.StatusInternalServerError,
	}, {
		name:       "add error, code already exists",
		input:      strings.NewReader(`{"url":"http://www.validUrl.com", "code":"d4jn8dsf"}`),
		statusCode: http.StatusConflict,
	}, {
		name:       "valid request",
		input:      strings.NewReader(`{"url":"http://www.validUrl.com"}`),
		statusCode: http.StatusCreated,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api", tc.input)
			rec := httptest.NewRecorder()

			c.Add(rec, req)

			result := rec.Result()

			if tc.statusCode != result.StatusCode {
				resBody, _ := ioutil.ReadAll(result.Body)
				t.Errorf("expected status code (%v), got (%v) with response: (%v)", tc.statusCode, result.StatusCode, string(resBody))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	s := ServiceMock{}
	l := log.New(os.Stdout, "urls-api", log.LstdFlags)
	c := NewController(&s, l)

	testCases := []struct {
		name       string
		input      string
		statusCode int
	}{
		{
			name:       "empty query",
			input:      "",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "non integer query",
			input:      "id",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "delete error",
			input:      "0",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "valid request",
			input:      "1",
			statusCode: http.StatusOK,
		},
		{
			name:       "negative id",
			input:      "-1",
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/"+tc.input, nil)
			rec := httptest.NewRecorder()

			c.Delete(rec, req)
			result := rec.Result()

			if result.StatusCode != tc.statusCode {
				resBody, _ := ioutil.ReadAll(result.Body)
				t.Errorf("expected status code (%v), got (%v) with response: (%v)", tc.statusCode, result.StatusCode, string(resBody))
			}
		})
	}
}

func TestGet(t *testing.T) {
	s := ServiceMock{}
	l := log.New(os.Stdout, "urls-api", log.LstdFlags)
	c := NewController(&s, l)

	testCases := []struct {
		name       string
		input      string
		statusCode int
	}{
		{
			name:       "empty query",
			input:      "",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "non integer query",
			input:      "id",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "get error",
			input:      "0",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "valid request",
			input:      "1",
			statusCode: http.StatusOK,
		},
		{
			name:       "valid request, url not found",
			input:      "-1",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/"+tc.input, nil)
			rec := httptest.NewRecorder()

			c.Get(rec, req)
			result := rec.Result()

			if result.StatusCode != tc.statusCode {
				resBody, _ := ioutil.ReadAll(result.Body)
				t.Errorf("expected status code (%v), got (%v) with response: (%v)", tc.statusCode, result.StatusCode, string(resBody))
			}
		})
	}
}

func TestRedirectShortUrl(t *testing.T) {
	s := ServiceMock{}
	l := log.New(os.Stdout, "urls-api", log.LstdFlags)
	c := NewController(&s, l)

	testCases := []struct {
		name       string
		input      string
		statusCode int
	}{
		{
			name:       "empty query",
			input:      "",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "get error",
			input:      "invalidCode",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "valid request",
			input:      "84gfj4i9",
			statusCode: http.StatusFound,
		},
		{
			name:       "valid request, url not found",
			input:      "84gfasdf",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/"+tc.input, nil)
			rec := httptest.NewRecorder()

			c.RedirectShortUrl(rec, req)
			result := rec.Result()

			if result.StatusCode != tc.statusCode {
				resBody, _ := ioutil.ReadAll(result.Body)
				t.Errorf("expected status code (%v), got (%v) with response: (%v)", tc.statusCode, result.StatusCode, string(resBody))
			}
		})
	}
}

func TestGetCounter(t *testing.T) {
	s := ServiceMock{}
	l := log.New(os.Stdout, "urls-api", log.LstdFlags)
	c := NewController(&s, l)

	testCases := []struct {
		name       string
		input      string
		statusCode int
	}{
		{
			name:       "empty query",
			input:      "",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "non integer query",
			input:      "id",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "get error",
			input:      "0",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "valid request",
			input:      "1",
			statusCode: http.StatusOK,
		},
		{
			name:       "valid request, url not found",
			input:      "-1",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/"+tc.input, nil)
			rec := httptest.NewRecorder()

			c.GetCounter(rec, req)
			result := rec.Result()

			if result.StatusCode != tc.statusCode {
				resBody, _ := ioutil.ReadAll(result.Body)
				t.Errorf("expected status code (%v), got (%v) with response: (%v)", tc.statusCode, result.StatusCode, string(resBody))
			}
		})
	}
}
