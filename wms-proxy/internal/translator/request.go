package translator

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"wms-proxy/internal/services"
	"wms-proxy/internal/transform"
	"wms-proxy/pkg/wms"
)

// TranslateWMSToArcGIS converts WMS GetMap parameters to ArcGIS REST export parameters
func TranslateWMSToArcGIS(wmsParams *wms.WMSParams) (*wms.ArcGISParams, error) {
	return TranslateWMSToArcGISWithTransform(wmsParams, nil)
}

// TranslateWMSToArcGISWithTransform converts WMS GetMap parameters to ArcGIS REST export parameters
// with optional coordinate transformation
func TranslateWMSToArcGISWithTransform(wmsParams *wms.WMSParams, transformer *transform.CoordinateTransformer) (*wms.ArcGISParams, error) {
	bbox := wmsParams.BBOX
	sourceSRS := wmsParams.GetSRS()
	targetSRS := translateSRS(sourceSRS)

	// If we have a transformer, check if transformation is needed
	if transformer != nil && sourceSRS != "" && targetSRS != "" {
		// The sourceSRS parameter indicates the coordinate system of the incoming bbox coordinates
		fromCRS := transformer.NormalizeCRS(sourceSRS)
		// The backend ArcGIS server expects coordinates in EPSG:3424 (New Jersey State Plane)
		toCRS := "EPSG:3424"

		// Only transform if the CRS are different
		if fromCRS != toCRS {
			transformedBBox, err := transformer.TransformBBox(bbox, fromCRS, toCRS)
			if err != nil {
				// Log the error but continue with original bbox
				// In a production system, you might want to handle this differently
				fmt.Printf("Warning: coordinate transformation failed: %v\n", err)
			} else {
				bbox = transformedBBox
			}
		}
	}

	arcgisParams := &wms.ArcGISParams{
		BBOX:        bbox,
		Size:        fmt.Sprintf("%d,%d", wmsParams.Width, wmsParams.Height),
		Format:      translateFormat(wmsParams.Format),
		BBoxSR:      targetSRS,
		ImageSR:     targetSRS,
		Layers:      translateLayers(wmsParams.Layers),
		Transparent: translateTransparent(wmsParams.Transparent),
		DPI:         96,
		F:           "image",
	}

	return arcgisParams, nil
}

// TranslateWMSToArcGISWithTransformAndBackendSR converts WMS GetMap parameters to ArcGIS REST export parameters
// with coordinate transformation using dynamic backend SR detection
func TranslateWMSToArcGISWithTransformAndBackendSR(wmsParams *wms.WMSParams, transformer *transform.CoordinateTransformer, srDetector *services.BackendSRDetector, ctx context.Context, servicePath string) (*wms.ArcGISParams, error) {
	bbox := wmsParams.BBOX
	sourceSRS := wmsParams.GetSRS()

	// If we have a transformer and SR detector, use dynamic backend SR detection
	if transformer != nil && srDetector != nil && sourceSRS != "" {
		// The sourceSRS parameter indicates the coordinate system of the incoming bbox coordinates
		fromCRS := transformer.NormalizeCRS(sourceSRS)

		// Detect what spatial reference system the backend service expects
		toCRS, err := srDetector.GetBackendSR(ctx, servicePath)
		if err != nil {
			// Log the error but continue with original bbox
			fmt.Printf("Warning: failed to detect backend SR: %v\n", err)
			toCRS = "EPSG:3424" // Fallback
		}

		// Only transform if the CRS are different
		if fromCRS != toCRS {
			transformedBBox, err := transformer.TransformBBox(bbox, fromCRS, toCRS)
			if err != nil {
				// Log the error but continue with original bbox
				fmt.Printf("Warning: coordinate transformation failed: %v\n", err)
			} else {
				bbox = transformedBBox
			}
		}
	}

	arcgisParams := &wms.ArcGISParams{
		BBOX:        bbox,
		Size:        fmt.Sprintf("%d,%d", wmsParams.Width, wmsParams.Height),
		Format:      translateFormat(wmsParams.Format),
		Transparent: translateTransparent(wmsParams.Transparent),
		Layers:      translateLayers(wmsParams.Layers),
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

	originalSRS := strings.TrimSpace(srs)
	upperSRS := strings.ToUpper(originalSRS)

	// Handle common formats
	if strings.HasPrefix(upperSRS, "EPSG:") {
		return upperSRS
	}

	// Handle some common variations
	switch upperSRS {
	case "3857", "900913":
		return "EPSG:3857"
	case "4326":
		return "EPSG:4326"
	case "3424":
		return "EPSG:3424"
	default:
		// If it's just a number, add EPSG: prefix
		if _, err := strconv.Atoi(originalSRS); err == nil {
			return "EPSG:" + originalSRS
		}
		// Return original case for unrecognized values
		return originalSRS
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
