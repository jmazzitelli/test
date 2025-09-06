package wms

import (
	"fmt"
	"strconv"
	"strings"
)

// WMSParams represents parsed WMS request parameters
type WMSParams struct {
	Service     string
	Version     string
	Request     string
	Layers      string
	Styles      string
	Format      string
	Transparent string
	BGColor     string
	SRS         string
	CRS         string
	BBOX        string
	Width       int
	Height      int
}

// ArcGISParams represents ArcGIS REST API parameters
type ArcGISParams struct {
	BBOX        string
	Size        string
	Format      string
	BBoxSR      string
	ImageSR     string
	Layers      string
	Transparent string
	DPI         int
	F           string
}

// ParseWMSParams extracts WMS parameters from query values
func ParseWMSParams(queryParams map[string][]string) (*WMSParams, error) {
	params := &WMSParams{}

	// Helper function to get first value from query params (case-insensitive)
	getValue := func(key string) string {
		// Try exact match first
		if values, exists := queryParams[key]; exists && len(values) > 0 {
			return values[0]
		}
		// Try case-insensitive match
		for k, v := range queryParams {
			if strings.EqualFold(k, key) && len(v) > 0 {
				return v[0]
			}
		}
		return ""
	}

	params.Service = getValue("SERVICE")
	params.Version = getValue("VERSION")
	params.Request = getValue("REQUEST")
	params.Layers = getValue("LAYERS")
	params.Styles = getValue("STYLES")
	params.Format = getValue("FORMAT")
	params.Transparent = getValue("TRANSPARENT")
	params.BGColor = getValue("BGCOLOR")
	params.SRS = getValue("SRS")
	params.CRS = getValue("CRS")
	params.BBOX = getValue("BBOX")

	// Parse width and height
	if widthStr := getValue("WIDTH"); widthStr != "" {
		if width, err := strconv.Atoi(widthStr); err == nil {
			params.Width = width
		} else {
			return nil, fmt.Errorf("invalid WIDTH parameter: %s", widthStr)
		}
	}

	if heightStr := getValue("HEIGHT"); heightStr != "" {
		if height, err := strconv.Atoi(heightStr); err == nil {
			params.Height = height
		} else {
			return nil, fmt.Errorf("invalid HEIGHT parameter: %s", heightStr)
		}
	}

	// Validate required parameters for GetMap
	if strings.EqualFold(params.Request, "GetMap") {
		if params.BBOX == "" {
			return nil, fmt.Errorf("BBOX parameter is required for GetMap requests")
		}
		if params.Width <= 0 || params.Height <= 0 {
			return nil, fmt.Errorf("WIDTH and HEIGHT parameters are required and must be positive for GetMap requests")
		}
		if params.Layers == "" {
			return nil, fmt.Errorf("LAYERS parameter is required for GetMap requests")
		}
	}

	return params, nil
}

// GetSRS returns the spatial reference system, preferring CRS over SRS
func (p *WMSParams) GetSRS() string {
	if p.CRS != "" {
		return p.CRS
	}
	return p.SRS
}
