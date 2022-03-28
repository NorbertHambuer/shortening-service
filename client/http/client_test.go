package http

import (
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"testing"
)

func TestCreate(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")

		var url Url
		err := url.FromJSON(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "unable to decode url"}`))
			return
		}

		if url.Url == "www.invalidUrl.com" {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "invalid url"}`))
			return
		}

		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(`{"id":6,"code":"VgUJPzDN","url":"https://www.google.ro/search?q=some","shortUrl":"http://localhost:3000/VgUJPzDN","domain":"http://localhost:3000","counter":2}`))
	}))

	client := NewClient(svr.URL)

	testCases := []struct {
		name    string
		input   CreateRequest
		isError bool
	}{
		{
			name:    "empty request",
			input:   CreateRequest{},
			isError: true,
		},
		{
			name:    "invalid url request",
			input:   CreateRequest{Url: "www.invalidUrl.com"},
			isError: true,
		},
		{
			name:    "valid request",
			input:   CreateRequest{Url: "www.validUrl.com"},
			isError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.Create(tc.input)

			if (err != nil) != tc.isError {
				t.Errorf("expected error (%v), got error (%v)", tc.isError, err)
			}
		})
	}
}

func TestDelete(t *testing.T){
	svr := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")

		id, err := strconv.Atoi(path.Base(r.URL.String()))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "invalid id value}`))
			return
		}

		if id == 0 {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "error deleting id"}`))
			return
		}

		rw.WriteHeader(http.StatusOK)
	}))

	client := NewClient(svr.URL)

	testCases := []struct {
		name    string
		input   int64
		isError bool
	}{
		{
			name:    "delete request error",
			input:   0,
			isError: true,
		},
		{
			name:    "valid request",
			input:   1,
			isError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := client.Delete(tc.input)

			if (err != nil) != tc.isError {
				t.Errorf("expected error (%v), got error (%v)", tc.isError, err)
			}
		})
	}
}

func TestGet(t *testing.T){
	svr := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")

		id, err := strconv.Atoi(path.Base(r.URL.String()))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "invalid id value}`))
			return
		}

		if id == -1 {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		if id == 0 {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "error deleting id"}`))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"id":6,"code":"VgUJPzDN","url":"https://www.google.ro/search?q=some","shortUrl":"http://localhost:3000/VgUJPzDN","domain":"http://localhost:3000","counter":2}`))
	}))

	client := NewClient(svr.URL)

	testCases := []struct {
		name    string
		input   int64
		isError bool
	}{
		{
			name:    "get request error",
			input:   0,
			isError: true,
		},
		{
			name:    "valid request",
			input:   1,
			isError: false,
		},
		{
			name:    "url not found",
			input:   -1,
			isError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.Get(tc.input)

			if (err != nil) != tc.isError {
				t.Errorf("expected error (%v), got error (%v)", tc.isError, err)
			}
		})
	}
}

func TestGetCounter(t *testing.T){
	svr := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")

		id, err := strconv.Atoi(path.Base(r.URL.String()))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "invalid id value}`))
			return
		}

		if id == -1 {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		if id == 0 {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(`{"message": "error deleting id"}`))
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"counter":6}`))
	}))

	client := NewClient(svr.URL)

	testCases := []struct {
		name    string
		input   int64
		isError bool
	}{
		{
			name:    "get request error",
			input:   0,
			isError: true,
		},
		{
			name:    "valid request",
			input:   1,
			isError: false,
		},
		{
			name:    "url not found",
			input:   -1,
			isError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.GetCounter(tc.input)

			if (err != nil) != tc.isError {
				t.Errorf("expected error (%v), got error (%v)", tc.isError, err)
			}
		})
	}
}
