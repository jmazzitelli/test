package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ArcGISClient handles HTTP requests to ArcGIS REST API
type ArcGISClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewArcGISClient creates a new ArcGIS REST API client
func NewArcGISClient(baseURL string, timeout time.Duration) *ArcGISClient {
	return &ArcGISClient{
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		baseURL: baseURL,
	}
}

// Get performs a GET request to the ArcGIS server
func (c *ArcGISClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set appropriate headers
	req.Header.Set("User-Agent", "WMS-Proxy/1.0")
	req.Header.Set("Accept", "image/png,image/jpeg,image/gif,*/*")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// HealthCheck verifies connectivity to the ArcGIS server
func (c *ArcGISClient) HealthCheck(ctx context.Context) error {
	// Try to access the base ArcGIS REST services endpoint
	healthURL := c.baseURL + "/arcgis/rest/services"

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}
