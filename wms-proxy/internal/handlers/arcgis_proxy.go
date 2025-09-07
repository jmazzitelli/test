package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"wms-proxy/internal/client"
	"wms-proxy/internal/services"
	"wms-proxy/internal/transform"
	"wms-proxy/internal/translator"
)

// ArcGISProxyHandler handles direct ArcGIS REST API requests and proxies them
type ArcGISProxyHandler struct {
	arcgisClient client.ArcGISClientInterface
	logger       *slog.Logger
	baseURL      string
	transformer  *transform.CoordinateTransformer
	srDetector   *services.BackendSRDetector
}

// NewArcGISProxyHandler creates a new ArcGIS proxy handler
func NewArcGISProxyHandler(arcgisClient client.ArcGISClientInterface, logger *slog.Logger, baseURL string) *ArcGISProxyHandler {
	return &ArcGISProxyHandler{
		arcgisClient: arcgisClient,
		logger:       logger,
		baseURL:      baseURL,
		transformer:  transform.NewCoordinateTransformer(),
		srDetector:   services.NewBackendSRDetector(arcgisClient, logger),
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

	// Parse and potentially transform coordinates in the query parameters
	targetURL, err := h.buildTransformedURL(r)
	if err != nil {
		h.logger.Error("Failed to build transformed URL", "error", err)
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
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

// buildTransformedURL builds the target URL with coordinate transformation if needed
func (h *ArcGISProxyHandler) buildTransformedURL(r *http.Request) (string, error) {
	// Parse the query parameters
	queryParams := r.URL.Query()

	// Check if we need to transform coordinates
	bbox := queryParams.Get("bbox")
	bboxSR := queryParams.Get("bboxSR")

	// If we have both bbox and bboxSR, check if transformation is needed
	if bbox != "" && bboxSR != "" {
		// The bboxSR parameter indicates the coordinate system of the incoming bbox coordinates
		fromCRS := h.transformer.NormalizeCRS(bboxSR)

		// Detect what spatial reference system the backend service expects
		toCRS, err := h.srDetector.GetBackendSR(r.Context(), r.URL.Path)
		if err != nil {
			h.logger.Warn("Failed to detect backend SR, using fallback",
				"error", err,
				"service_path", r.URL.Path,
				"fallback_sr", "EPSG:3424")
			toCRS = "EPSG:3424" // Fallback to New Jersey State Plane
		}

		// Only transform if the CRS are different
		if fromCRS != toCRS {
			// Transform the coordinates from the client's CRS to what the backend expects
			transformedBBox, err := h.transformer.TransformBBox(bbox, fromCRS, toCRS)
			if err != nil {
				h.logger.Warn("Coordinate transformation failed, using original bbox",
					"error", err,
					"original_bbox", bbox,
					"from_crs", fromCRS,
					"to_crs", toCRS,
				)
				// Continue with original bbox if transformation fails
			} else {
				// Use the transformed coordinates and update bboxSR to match backend expectation
				queryParams.Set("bbox", transformedBBox)
				// Extract just the numeric part for bboxSR (e.g., "EPSG:3424" -> "3424")
				if strings.HasPrefix(toCRS, "EPSG:") {
					queryParams.Set("bboxSR", strings.TrimPrefix(toCRS, "EPSG:"))
				} else {
					queryParams.Set("bboxSR", toCRS)
				}
				h.logger.Info("Transformed coordinates",
					"original_bbox", bbox,
					"transformed_bbox", transformedBBox,
					"from_crs", fromCRS,
					"to_crs", toCRS,
				)
			}
		}
	}

	// Build the target URL
	targetURL := h.baseURL + r.URL.Path
	if len(queryParams) > 0 {
		targetURL += "?" + queryParams.Encode()
	}

	return targetURL, nil
}

// isArcGISPath checks if the request path is for ArcGIS REST API
func isArcGISPath(path string) bool {
	return strings.HasPrefix(path, "/arcgis/")
}
