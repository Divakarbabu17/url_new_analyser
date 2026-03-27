package linkchecker

import (
	"net/http"
	"sync"
	"time"

	"url_new_analyser/internal/core/ports"
)

// WorkerPool implements ports.LinkChecker interface
type WorkerPool struct {
	Client      *http.Client
	concurrency int
	stopCh      chan struct{}
}

// NewWorkerPoolLinkChecker creates a link checker worker pool
func NewWorkerPoolLinkChecker(concurrency int, timeout time.Duration) *WorkerPool {
	return &WorkerPool{
		Client: &http.Client{
			Timeout: timeout,
		},
		concurrency: concurrency,
		stopCh:      make(chan struct{}),
	}
}

// CheckLinks checks multiple URLs concurrently
func (w *WorkerPool) CheckLinks(urls []string) []ports.LinkResult {
	results := make([]ports.LinkResult, len(urls))
	var wg sync.WaitGroup
	sem := make(chan struct{}, w.concurrency)

	for i, link := range urls {
		wg.Add(1)
		go func(i int, link string) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-w.stopCh:
				results[i] = ports.LinkResult{URL: link, StatusCode: 0, OK: false}
				return
			}

			resp, err := w.Client.Head(link)
			if err != nil {
				results[i] = ports.LinkResult{URL: link, StatusCode: 0, OK: false}
				return
			}
			resp.Body.Close()
			results[i] = ports.LinkResult{
				URL:        link,
				StatusCode: resp.StatusCode,
				OK:         resp.StatusCode < 400,
			}
		}(i, link)
	}

	wg.Wait()
	return results
}

// Stop gracefully stops all workers
func (w *WorkerPool) Stop() {
	close(w.stopCh)
}

// ✅ Compile-time interface assertion
var _ ports.LinkChecker = (*WorkerPool)(nil)