package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Create calls the POST /api endpoint of the shortening service url that adds a new short url to the database
func (c *Client) Create(r CreateRequest) (Url, error) {
	// validate request
	if err := r.Validate(); err != nil {
		return Url{}, err
	}

	buf := new(bytes.Buffer)
	err := r.ToJSON(buf)
	if err != nil {
		return Url{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api", c.BaseURL), buf)
	if err != nil {
		return Url{}, err
	}

	req.Header.Add("Content-type", "application/json")

	// call endpoint
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Url{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		var errMsg ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errMsg)
		if err != nil {
			return Url{}, err
		}

		return Url{}, fmt.Errorf("error calling the create endpoint: %s", errMsg)
	}

	// decode response
	var u Url
	err = u.FromJSON(resp.Body)
	if err != nil {
		return Url{}, err
	}

	return u, nil
}

// Delete calls the DELETE /api endpoint of the shortening service url that deletes the url with the given ID
func (c *Client) Delete(id int64) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/%d", c.BaseURL, id), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var errMsg ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errMsg)
		if err != nil {
			return err
		}

		return fmt.Errorf("error calling the delete endpoint: %s", errMsg)
	}

	return nil
}

// Get calls the GET /api endpoint of the shortening service url that returns the url object with the given ID
func (c *Client) Get(id int64) (Url, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%d", c.BaseURL, id), nil)
	if err != nil {
		return Url{}, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return Url{}, err
	}

	if resp.StatusCode == http.StatusNotFound{
		return Url{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		var errMsg ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errMsg)
		if err != nil {
			return Url{}, err
		}

		return Url{}, fmt.Errorf("error calling the delete endpoint: %s", errMsg)
	}

	// decode response
	var u Url
	err = u.FromJSON(resp.Body)
	if err != nil {
		return Url{}, err
	}

	return u, nil
}

// GetCounter calls the GET /counter endpoint of the shortening service url that returns the number of redirections for the given id
func (c *Client) GetCounter(id int64) (int64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/%d", c.BaseURL, id), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusNotFound{
		return 0, nil
	}

	if resp.StatusCode != http.StatusOK {
		var errMsg ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errMsg)
		if err != nil {
			return 0, err
		}

		return 0, fmt.Errorf("error calling the delete endpoint: %s", errMsg)
	}

	var cr CounterResponse
	err = cr.FromJSON(resp.Body)
	if err != nil {
		return 0, err
	}

	return cr.Value, nil
}

// NewClient creates a new client object
func NewClient(baseUrl string) *Client {
	return &Client{
		BaseURL: baseUrl,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}
