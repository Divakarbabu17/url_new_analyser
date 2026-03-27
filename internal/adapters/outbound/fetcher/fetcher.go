package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"url_new_analyser/internal/core/ports"
)

// HTTPFetcher implements the Fetcher interface
type HTTPFetcher struct {
	Client *http.Client
}

// NewHTTPFetcher creates a new HTTPFetcher with timeout
func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	return &HTTPFetcher{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Fetch performs GET request to the given URL
func (f *HTTPFetcher) Fetch(url string) (string, int, error) {
	resp, err := f.Client.Get(url)
	if err != nil {
		return "", 0, fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, fmt.Errorf("failed to read body: %w", err)
	}

	return string(body), resp.StatusCode, nil
}

// Ensure it implements Fetcher port
var _ ports.Fetcher = (*HTTPFetcher)(nil)
