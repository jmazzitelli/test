package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ServiceMetadata represents ArcGIS MapServer service metadata
type ServiceMetadata struct {
	SpatialReference struct {
		WKID       int    `json:"wkid"`
		LatestWKID int    `json:"latestWkid"`
		WKT        string `json:"wkt"`
	} `json:"spatialReference"`
	SupportedQueryFormats interface{} `json:"supportedQueryFormats"` // Can be string or []string
	MaxRecordCount        int         `json:"maxRecordCount"`
	Capabilities          string      `json:"capabilities"`
}

// ArcGISClientInterface defines the interface for ArcGIS client operations
type ArcGISClientInterface interface {
	Get(ctx context.Context, url string) (*http.Response, error)
	GetServiceMetadata(ctx context.Context, servicePath string) (*ServiceMetadata, error)
}

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

// GetServiceMetadata retrieves metadata for an ArcGIS MapServer service
func (c *ArcGISClient) GetServiceMetadata(ctx context.Context, servicePath string) (*ServiceMetadata, error) {
	// Remove /export suffix if present to get the service root
	serviceRoot := strings.TrimSuffix(servicePath, "/export")

	// Build metadata URL
	metadataURL := c.baseURL + serviceRoot + "?f=json"

	req, err := http.NewRequestWithContext(ctx, "GET", metadataURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create metadata request: %w", err)
	}

	req.Header.Set("User-Agent", "WMS-Proxy/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute metadata request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("metadata request failed with status: %d", resp.StatusCode)
	}

	var metadata ServiceMetadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode metadata response: %w", err)
	}

	return &metadata, nil
}
