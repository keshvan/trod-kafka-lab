package dataservice

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type DataClient struct {
	baseURL string
	client  *http.Client
}

func NewDataClient(baseURL string) *DataClient {
	return &DataClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *DataClient) Get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("build data service url: %w", err)
	}

	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request to data service: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call data service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read data service response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("data service returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
