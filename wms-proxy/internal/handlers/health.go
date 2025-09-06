package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"wms-proxy/internal/client"
)

// HealthHandler provides health check functionality
type HealthHandler struct {
	arcgisClient *client.ArcGISClient
}

// NewHealthHandler creates a new health check handler
func NewHealthHandler(arcgisClient *client.ArcGISClient) *HealthHandler {
	return &HealthHandler{
		arcgisClient: arcgisClient,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Upstream  string `json:"upstream"`
	Message   string `json:"message,omitempty"`
}

// ServeHTTP handles health check requests
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Upstream:  "ok",
	}

	// Check upstream ArcGIS server connectivity
	if err := h.arcgisClient.HealthCheck(ctx); err != nil {
		response.Status = "unhealthy"
		response.Upstream = "error"
		response.Message = err.Error()
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode health response", http.StatusInternalServerError)
		return
	}
}
