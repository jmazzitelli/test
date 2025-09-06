package handlers

import (
	"net/http"

	"wms-proxy/pkg/wms"
)

// CapabilitiesHandler handles WMS GetCapabilities requests
type CapabilitiesHandler struct {
	baseURL string
}

// NewCapabilitiesHandler creates a new capabilities handler
func NewCapabilitiesHandler(baseURL string) *CapabilitiesHandler {
	return &CapabilitiesHandler{
		baseURL: baseURL,
	}
}

// ServeHTTP handles GetCapabilities requests
func (h *CapabilitiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Generate capabilities XML
	capabilitiesXML, err := wms.GenerateCapabilities(h.baseURL)
	if err != nil {
		http.Error(w, "Failed to generate capabilities", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", "application/vnd.ogc.wms_xml")
	w.Header().Set("Cache-Control", "max-age=3600") // Cache for 1 hour

	w.WriteHeader(http.StatusOK)
	w.Write(capabilitiesXML)
}
