package handlers

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"wms-proxy/internal/client"
)

// mockArcGISClient is a mock implementation of ArcGISClientInterface for testing
type mockArcGISClient struct {
	lastRequestURL string
	response       *http.Response
	err            error
}

// Ensure mockArcGISClient implements client.ArcGISClientInterface
var _ client.ArcGISClientInterface = (*mockArcGISClient)(nil)

func (m *mockArcGISClient) Get(ctx context.Context, url string) (*http.Response, error) {
	m.lastRequestURL = url
	return m.response, m.err
}

func (m *mockArcGISClient) GetServiceMetadata(ctx context.Context, servicePath string) (*client.ServiceMetadata, error) {
	// Return a mock metadata for testing
	return &client.ServiceMetadata{
		SpatialReference: struct {
			WKID       int    `json:"wkid"`
			LatestWKID int    `json:"latestWkid"`
			WKT        string `json:"wkt"`
		}{
			WKID: 3424,
		},
	}, nil
}

func TestArcGISProxyHandler_buildTransformedURL(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	mockClient := &mockArcGISClient{}
	handler := NewArcGISProxyHandler(mockClient, logger, "https://example.com")

	tests := []struct {
		name            string
		requestURL      string
		expectedBBox    string
		shouldTransform bool
	}{
		{
			name:            "No transformation needed - client sends EPSG:3424, backend expects EPSG:3424",
			requestURL:      "/arcgis/rest/services/test/MapServer/export?bbox=629066,684288,629793,685020&bboxSR=3424&imageSR=3424&size=256,256&f=image",
			shouldTransform: false,
		},
		{
			name:            "Transform EPSG:4326 to EPSG:3424 - client sends EPSG:4326, backend expects EPSG:3424",
			requestURL:      "/arcgis/rest/services/test/MapServer/export?bbox=-74.006000,40.710974,-74.003364,40.712972&bboxSR=4326&imageSR=4326&size=256,256&f=image",
			shouldTransform: true,
		},
		{
			name:            "Transform EPSG:3857 to EPSG:3424 - client sends EPSG:3857, backend expects EPSG:3424",
			requestURL:      "/arcgis/rest/services/test/MapServer/export?bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=3857&imageSR=3857&size=256,256&f=image",
			shouldTransform: true,
		},
		{
			name:            "No bbox parameter",
			requestURL:      "/arcgis/rest/services/test/MapServer/export?size=256,256&f=image",
			shouldTransform: false,
		},
		{
			name:            "No bboxSR parameter",
			requestURL:      "/arcgis/rest/services/test/MapServer/export?bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&size=256,256&f=image",
			shouldTransform: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a mock request
			req := httptest.NewRequest("GET", test.requestURL, nil)

			// Call buildTransformedURL
			targetURL, err := handler.buildTransformedURL(req)
			if err != nil {
				t.Errorf("buildTransformedURL failed: %v", err)
				return
			}

			// Parse the result URL to check parameters
			parsedURL, err := url.Parse(targetURL)
			if err != nil {
				t.Errorf("failed to parse result URL: %v", err)
				return
			}

			originalBBox := req.URL.Query().Get("bbox")
			resultBBox := parsedURL.Query().Get("bbox")

			if test.shouldTransform {
				if originalBBox != "" && resultBBox == originalBBox {
					t.Errorf("expected bbox to be transformed, but it remained the same: %s", resultBBox)
				}
				if resultBBox == "" && originalBBox != "" {
					t.Error("bbox was removed during transformation")
				}
			} else {
				if originalBBox != "" && resultBBox != originalBBox {
					t.Errorf("bbox was unexpectedly transformed: original=%s, result=%s", originalBBox, resultBBox)
				}
			}

			// Verify the base URL is correct
			expectedBase := "https://example.com" + req.URL.Path
			if !strings.HasPrefix(targetURL, expectedBase) {
				t.Errorf("target URL doesn't have expected base: got %s, expected to start with %s", targetURL, expectedBase)
			}

			t.Logf("Original bbox: %s", originalBBox)
			t.Logf("Transformed bbox: %s", resultBBox)
			t.Logf("Target URL: %s", targetURL)
		})
	}
}

func TestArcGISProxyHandler_ServeHTTP(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name           string
		method         string
		requestURL     string
		mockResponse   *http.Response
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Successful GET request with coordinate transformation",
			method:         "GET",
			requestURL:     "/arcgis/rest/services/test/MapServer/export?bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=3424&size=256,256&f=image",
			mockResponse:   &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("mock image data")), Header: make(http.Header)},
			expectedStatus: 200,
		},
		{
			name:           "POST method not allowed",
			method:         "POST",
			requestURL:     "/arcgis/rest/services/test/MapServer/export",
			expectedStatus: 405,
		},
		{
			name:           "Successful GET request without transformation",
			method:         "GET",
			requestURL:     "/arcgis/rest/services/test/MapServer/export?size=256,256&f=image",
			mockResponse:   &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("mock image data")), Header: make(http.Header)},
			expectedStatus: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := &mockArcGISClient{
				response: test.mockResponse,
				err:      test.mockError,
			}
			handler := NewArcGISProxyHandler(mockClient, logger, "https://example.com")

			req := httptest.NewRequest(test.method, test.requestURL, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != test.expectedStatus {
				t.Errorf("expected status %d, got %d", test.expectedStatus, w.Code)
			}

			// For successful requests, verify that the mock client was called
			if test.expectedStatus == 200 && mockClient.lastRequestURL == "" {
				t.Error("expected mock client to be called, but it wasn't")
			}

			// For transformation cases, verify the URL was modified
			if test.expectedStatus == 200 && strings.Contains(test.requestURL, "bboxSR=3424") {
				if !strings.Contains(mockClient.lastRequestURL, "bbox=") {
					t.Error("expected transformed bbox in the request to upstream server")
				}
				t.Logf("Upstream request URL: %s", mockClient.lastRequestURL)
			}
		})
	}
}

func TestArcGISProxyHandler_CoordinateTransformationIntegration(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	// Create a mock response
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("mock image data")),
		Header:     make(http.Header),
	}
	mockResponse.Header.Set("Content-Type", "image/png")

	mockClient := &mockArcGISClient{
		response: mockResponse,
	}
	handler := NewArcGISProxyHandler(mockClient, logger, "https://mapserver.example.com")

	// Test request with coordinate transformation
	// Client sends EPSG:3857 coordinates, backend expects EPSG:3424 (from mock)
	originalBBox := "-8238310.24,4969803.4,-8238016.75,4970096.9"
	requestURL := "/arcgis/rest/services/Features/Environmental_admin/MapServer/export?" +
		"bbox=" + originalBBox +
		"&bboxSR=3857" + // Client sends Web Mercator coordinates
		"&imageSR=3424" +
		"&size=256,256" +
		"&f=image" +
		"&layers=show:17"

	req := httptest.NewRequest("GET", requestURL, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify the response
	if w.Code != 200 {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify that the upstream request was made with transformed coordinates
	if mockClient.lastRequestURL == "" {
		t.Fatal("mock client was not called")
	}

	upstreamURL, err := url.Parse(mockClient.lastRequestURL)
	if err != nil {
		t.Fatalf("failed to parse upstream URL: %v", err)
	}

	transformedBBox := upstreamURL.Query().Get("bbox")
	if transformedBBox == "" {
		t.Error("bbox parameter missing in upstream request")
	}

	if transformedBBox == originalBBox {
		t.Error("bbox was not transformed - coordinates are identical to original")
	}

	// Verify other parameters are preserved
	if upstreamURL.Query().Get("bboxSR") != "3424" {
		t.Error("bboxSR parameter not preserved")
	}

	if upstreamURL.Query().Get("size") != "256,256" {
		t.Error("size parameter not preserved")
	}

	t.Logf("Original bbox: %s", originalBBox)
	t.Logf("Transformed bbox: %s", transformedBBox)
	t.Logf("Full upstream URL: %s", mockClient.lastRequestURL)
}

// Test error handling in coordinate transformation
func TestArcGISProxyHandler_TransformationErrorHandling(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("mock image data")),
		Header:     make(http.Header),
	}

	mockClient := &mockArcGISClient{
		response: mockResponse,
	}
	handler := NewArcGISProxyHandler(mockClient, logger, "https://example.com")

	// Test with invalid bbox format - should not cause server error
	requestURL := "/arcgis/rest/services/test/MapServer/export?" +
		"bbox=invalid,bbox,format" +
		"&bboxSR=3424" +
		"&size=256,256" +
		"&f=image"

	req := httptest.NewRequest("GET", requestURL, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Should handle the error gracefully and continue with original bbox
	if w.Code != 200 {
		t.Errorf("expected status 200 (graceful error handling), got %d", w.Code)
	}

	// Verify that a request was still made to upstream (with original bbox)
	if mockClient.lastRequestURL == "" {
		t.Error("expected request to be made to upstream server despite transformation error")
	}
}
