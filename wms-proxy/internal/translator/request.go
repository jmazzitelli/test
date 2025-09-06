package translator

import (
	"fmt"
	"net/url"
	"strings"

	"wms-proxy/pkg/wms"
)

// TranslateWMSToArcGIS converts WMS GetMap parameters to ArcGIS REST export parameters
func TranslateWMSToArcGIS(wmsParams *wms.WMSParams) (*wms.ArcGISParams, error) {
	arcgisParams := &wms.ArcGISParams{
		BBOX:        wmsParams.BBOX,
		Size:        fmt.Sprintf("%d,%d", wmsParams.Width, wmsParams.Height),
		Format:      translateFormat(wmsParams.Format),
		BBoxSR:      translateSRS(wmsParams.GetSRS()),
		ImageSR:     translateSRS(wmsParams.GetSRS()),
		Layers:      translateLayers(wmsParams.Layers),
		Transparent: translateTransparent(wmsParams.Transparent),
		DPI:         96,
		F:           "image",
	}

	return arcgisParams, nil
}

// BuildArcGISURL constructs the full ArcGIS REST API URL
func BuildArcGISURL(baseURL, servicePath string, params *wms.ArcGISParams) string {
	u, err := url.Parse(baseURL + servicePath)
	if err != nil {
		// Fallback to simple concatenation if URL parsing fails
		return baseURL + servicePath + "?" + buildQueryString(params)
	}

	query := u.Query()
	query.Set("bbox", params.BBOX)
	query.Set("size", params.Size)
	query.Set("format", params.Format)
	query.Set("bboxSR", params.BBoxSR)
	query.Set("imageSR", params.ImageSR)
	query.Set("layers", params.Layers)
	query.Set("transparent", params.Transparent)
	query.Set("dpi", fmt.Sprintf("%d", params.DPI))
	query.Set("f", params.F)

	u.RawQuery = query.Encode()
	return u.String()
}

func buildQueryString(params *wms.ArcGISParams) string {
	values := url.Values{}
	values.Set("bbox", params.BBOX)
	values.Set("size", params.Size)
	values.Set("format", params.Format)
	values.Set("bboxSR", params.BBoxSR)
	values.Set("imageSR", params.ImageSR)
	values.Set("layers", params.Layers)
	values.Set("transparent", params.Transparent)
	values.Set("dpi", fmt.Sprintf("%d", params.DPI))
	values.Set("f", params.F)
	return values.Encode()
}

// translateFormat converts WMS format to ArcGIS format
func translateFormat(wmsFormat string) string {
	switch strings.ToLower(wmsFormat) {
	case "image/png", "png":
		return "png32"
	case "image/jpeg", "jpeg", "jpg":
		return "jpg"
	case "image/gif", "gif":
		return "gif"
	default:
		return "png32" // Default to PNG with transparency
	}
}

// translateSRS converts WMS SRS/CRS to ArcGIS spatial reference
func translateSRS(srs string) string {
	if srs == "" {
		return "EPSG:3857" // Default to Web Mercator
	}

	// Handle common formats
	srs = strings.ToUpper(srs)
	if strings.HasPrefix(srs, "EPSG:") {
		return srs
	}

	// Handle some common variations
	switch srs {
	case "3857", "900913":
		return "EPSG:3857"
	case "4326":
		return "EPSG:4326"
	default:
		return srs // Pass through as-is
	}
}

// translateLayers converts WMS layers to ArcGIS layers format
func translateLayers(wmsLayers string) string {
	if wmsLayers == "" {
		return "show:0" // Default to first layer
	}

	// For simplicity, assume single layer ID and convert to show format
	// In a real implementation, this might need more sophisticated parsing
	layers := strings.Split(wmsLayers, ",")
	var arcgisLayers []string

	for _, layer := range layers {
		layer = strings.TrimSpace(layer)
		if layer != "" {
			arcgisLayers = append(arcgisLayers, "show:"+layer)
		}
	}

	if len(arcgisLayers) == 0 {
		return "show:0"
	}

	return strings.Join(arcgisLayers, ",")
}

// translateTransparent converts WMS transparent parameter
func translateTransparent(transparent string) string {
	switch strings.ToLower(transparent) {
	case "true", "1", "yes":
		return "true"
	case "false", "0", "no":
		return "false"
	default:
		return "true" // Default to transparent
	}
}
