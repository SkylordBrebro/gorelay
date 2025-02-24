package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client provides HTTP functionality
type Client struct {
	client  *http.Client
	baseURL string
}

// NewClient creates a new HTTP client
func NewClient(baseURL string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// Get performs a GET request
func (c *Client) Get(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET request returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// Post performs a POST request
func (c *Client) Post(path string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %v", err)
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("POST request returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// GetJSON performs a GET request and unmarshals the response into v
func (c *Client) GetJSON(path string, v interface{}) error {
	data, err := c.Get(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// PostJSON performs a POST request and unmarshals the response into v
func (c *Client) PostJSON(path string, data interface{}, v interface{}) error {
	respData, err := c.Post(path, data)
	if err != nil {
		return err
	}

	return json.Unmarshal(respData, v)
}

// SetTimeout sets the client timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.Timeout = timeout
}

// SetBaseURL sets the base URL for requests
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}
