package urlhandler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// New is URLHandler constructor.
func New(fileName string, reqTimeout time.Duration) *URLHandler {
	return &URLHandler{
		fileName:   fileName,
		reqTimeout: reqTimeout,
	}
}

// URLHandler is url handler instance.
type URLHandler struct {
	fileName   string
	reqTimeout time.Duration
}

// Run runs url handler.
func (h *URLHandler) Run(ctx context.Context) error {
	f, err := os.ReadFile(h.fileName)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	urls := strings.Split(string(f), "\n")

	h.handler(ctx, urls)

	return nil
}

// CheckURL checks URL and make request to rawURL.
func (h *URLHandler) CheckURL(ctx context.Context, rawURL string) (Response, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return Response{}, fmt.Errorf("parse %q url: %w", rawURL, err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, h.reqTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return Response{}, fmt.Errorf("make GET request with context to %q: %w", u, err)
	}

	startTime := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("do request to %q: %w", u, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	responseTime := time.Since(startTime)

	return Response{
		URL:           u.String(),
		ContentLength: resp.ContentLength,
		ResponseTime:  responseTime,
	}, nil
}

func (h *URLHandler) handler(ctx context.Context, urls []string) {
	respCh := make(chan Response)
	defer close(respCh)

	errCh := make(chan error)
	defer close(errCh)

	go h.writer(ctx, respCh, errCh)

	wg := new(sync.WaitGroup)

	for _, rawURL := range urls {
		rawURL := rawURL
		wg.Add(1)

		go func() {
			defer wg.Done()

			resp, err := h.CheckURL(ctx, rawURL)
			if err == nil { // if NO error
				respCh <- resp
			} else {
				errCh <- err
			}
		}()
	}

	wg.Wait()
}

func (h *URLHandler) writer(ctx context.Context, respCh <-chan Response, errCh <-chan error) {
	for {
		select {
		case resp, ok := <-respCh:
			if !ok {
				return
			}

			fmt.Printf("url: %s, content length: %d, response time: %v\n", resp.URL, resp.ContentLength, resp.ResponseTime)

		case err, ok := <-errCh:
			if !ok {
				return
			}
			fmt.Printf("error: %v\n", err)

		case <-ctx.Done():
			return

		default:
			continue
		}
	}
}

// Response is instance for sucessfully response.
type Response struct {
	URL           string
	ContentLength int64
	ResponseTime  time.Duration
}
