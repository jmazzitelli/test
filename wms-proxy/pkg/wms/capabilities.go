package wms

import (
	"encoding/xml"
	"fmt"
)

// WMSCapabilities represents a basic WMS GetCapabilities response
type WMSCapabilities struct {
	XMLName xml.Name `xml:"WMS_Capabilities"`
	Version string   `xml:"version,attr"`
	Service Service  `xml:"Service"`
}

// Service represents the WMS service information
type Service struct {
	Name           string         `xml:"Name"`
	Title          string         `xml:"Title"`
	OnlineResource OnlineResource `xml:"OnlineResource"`
}

// OnlineResource represents a URL reference
type OnlineResource struct {
	XMLName xml.Name `xml:"OnlineResource"`
	Href    string   `xml:"xlink:href,attr"`
	XLink   string   `xml:"xmlns:xlink,attr"`
}

// GenerateCapabilities creates a basic WMS capabilities XML response
func GenerateCapabilities(baseURL string) ([]byte, error) {
	capabilities := WMSCapabilities{
		Version: "1.1.1",
		Service: Service{
			Name:  "WMS",
			Title: "ArcGIS REST to WMS Proxy",
			OnlineResource: OnlineResource{
				Href:  baseURL,
				XLink: "http://www.w3.org/1999/xlink",
			},
		},
	}

	// Add XML header
	xmlData, err := xml.MarshalIndent(capabilities, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal capabilities XML: %w", err)
	}

	// Prepend XML declaration
	result := []byte(xml.Header + string(xmlData))
	return result, nil
}
