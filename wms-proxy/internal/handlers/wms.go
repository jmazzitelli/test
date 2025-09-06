package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"wms-proxy/internal/client"
	"wms-proxy/internal/translator"
	"wms-proxy/pkg/wms"
)

// WMSHandler handles WMS requests and proxies them to ArcGIS REST API
type WMSHandler struct {
	arcgisClient *client.ArcGISClient
	logger       *slog.Logger
	baseURL      string
	servicePath  string
}

// NewWMSHandler creates a new WMS handler
func NewWMSHandler(arcgisClient *client.ArcGISClient, logger *slog.Logger, baseURL, servicePath string) *WMSHandler {
	return &WMSHandler{
		arcgisClient: arcgisClient,
		logger:       logger,
		baseURL:      baseURL,
		servicePath:  servicePath,
	}
}

// ServeHTTP handles WMS requests
func (h *WMSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Log incoming request
	h.logger.Info("Incoming WMS request",
		"method", r.Method,
		"url", r.URL.String(),
		"remote_addr", r.RemoteAddr,
	)

	// Only handle GET requests
	if r.Method != http.MethodGet {
		translator.GenerateWMSError(w, "Only GET method is supported", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	queryParams := r.URL.Query()

	// Convert to map[string][]string for our parser
	params := make(map[string][]string)
	for key, values := range queryParams {
		params[key] = values
	}

	// Parse WMS parameters
	wmsParams, err := wms.ParseWMSParams(params)
	if err != nil {
		h.logger.Error("Failed to parse WMS parameters", "error", err)
		translator.GenerateWMSError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Handle different WMS requests
	switch strings.ToUpper(wmsParams.Request) {
	case "GETCAPABILITIES":
		h.handleGetCapabilities(w, r)
	case "GETMAP":
		h.handleGetMap(w, r, wmsParams)
	default:
		translator.GenerateWMSError(w, "Unsupported request type: "+wmsParams.Request, http.StatusBadRequest)
	}

	// Log request completion
	duration := time.Since(startTime)
	h.logger.Info("Request completed",
		"duration_ms", duration.Milliseconds(),
		"request_type", wmsParams.Request,
	)
}

// handleGetCapabilities processes WMS GetCapabilities requests
func (h *WMSHandler) handleGetCapabilities(w http.ResponseWriter, r *http.Request) {
	capabilitiesXML, err := wms.GenerateCapabilities(h.baseURL)
	if err != nil {
		h.logger.Error("Failed to generate capabilities", "error", err)
		translator.GenerateWMSError(w, "Failed to generate capabilities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.ogc.wms_xml")
	w.Header().Set("Cache-Control", "max-age=3600")
	w.WriteHeader(http.StatusOK)
	w.Write(capabilitiesXML)
}

// handleGetMap processes WMS GetMap requests
func (h *WMSHandler) handleGetMap(w http.ResponseWriter, r *http.Request, wmsParams *wms.WMSParams) {
	// Translate WMS parameters to ArcGIS parameters
	arcgisParams, err := translator.TranslateWMSToArcGIS(wmsParams)
	if err != nil {
		h.logger.Error("Failed to translate WMS parameters", "error", err)
		translator.GenerateWMSError(w, "Parameter translation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Build ArcGIS URL
	arcgisURL := translator.BuildArcGISURL(h.baseURL, h.servicePath, arcgisParams)

	h.logger.Info("Proxying to ArcGIS",
		"arcgis_url", arcgisURL,
		"bbox", arcgisParams.BBOX,
		"size", arcgisParams.Size,
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Make request to ArcGIS server
	arcgisResp, err := h.arcgisClient.Get(ctx, arcgisURL)
	if err != nil {
		h.logger.Error("Failed to request from ArcGIS server", "error", err)
		translator.GenerateWMSError(w, "Upstream server error", http.StatusBadGateway)
		return
	}

	// Translate and return response
	if err := translator.TranslateArcGISResponse(arcgisResp, w); err != nil {
		h.logger.Error("Failed to translate ArcGIS response", "error", err)
		// Response may have already been written, so we can't send another error
		return
	}

	h.logger.Info("Successfully proxied response",
		"status_code", arcgisResp.StatusCode,
		"content_type", arcgisResp.Header.Get("Content-Type"),
	)
}
