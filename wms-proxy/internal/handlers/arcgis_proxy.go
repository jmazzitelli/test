package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"wms-proxy/internal/client"
	"wms-proxy/internal/translator"
)

// ArcGISProxyHandler handles direct ArcGIS REST API requests and proxies them
type ArcGISProxyHandler struct {
	arcgisClient *client.ArcGISClient
	logger       *slog.Logger
	baseURL      string
}

// NewArcGISProxyHandler creates a new ArcGIS proxy handler
func NewArcGISProxyHandler(arcgisClient *client.ArcGISClient, logger *slog.Logger, baseURL string) *ArcGISProxyHandler {
	return &ArcGISProxyHandler{
		arcgisClient: arcgisClient,
		logger:       logger,
		baseURL:      baseURL,
	}
}

// ServeHTTP handles ArcGIS REST API requests and proxies them directly
func (h *ArcGISProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Log incoming request
	h.logger.Info("Incoming ArcGIS proxy request",
		"method", r.Method,
		"path", r.URL.Path,
		"query", r.URL.RawQuery,
		"remote_addr", r.RemoteAddr,
	)

	// Only handle GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Build the target ArcGIS URL by combining base URL with the request path and query
	targetURL := h.baseURL + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	h.logger.Info("Proxying to ArcGIS",
		"target_url", targetURL,
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Make request to ArcGIS server
	arcgisResp, err := h.arcgisClient.Get(ctx, targetURL)
	if err != nil {
		h.logger.Error("Failed to request from ArcGIS server", "error", err)
		http.Error(w, "Upstream server error", http.StatusBadGateway)
		return
	}

	// Copy the response directly (no translation needed for direct proxy)
	if err := translator.TranslateArcGISResponse(arcgisResp, w); err != nil {
		h.logger.Error("Failed to copy ArcGIS response", "error", err)
		return
	}

	// Log request completion
	duration := time.Since(startTime)
	h.logger.Info("ArcGIS proxy request completed",
		"duration_ms", duration.Milliseconds(),
		"status_code", arcgisResp.StatusCode,
		"content_type", arcgisResp.Header.Get("Content-Type"),
	)
}

// isArcGISPath checks if the request path is for ArcGIS REST API
func isArcGISPath(path string) bool {
	return strings.HasPrefix(path, "/arcgis/")
}
