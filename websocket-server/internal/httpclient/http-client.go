package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	Client *http.Client
}

func NewHttpClient(timeout time.Duration) *HttpClient {
	if timeout == 0 {
		return &HttpClient{
			Client: &http.Client{},
		}
	}
	return &HttpClient{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Generic method to make GET requests
func (h *HttpClient) Get(url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	h.setHeaders(req, headers)
	return h.Client.Do(req)
}

// Generic method to make POST requests
func (h *HttpClient) Post(url string, body interface{}, headers map[string]string) (*http.Response, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	h.setHeaders(req, headers)
	return h.Client.Do(req)
}

// Utility method to set headers on a request
func (h *HttpClient) setHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

// Helper method to parse the response body
func ParseResponseBody(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, result)
}
