package translator

import (
	"strings"
	"testing"

	"wms-proxy/internal/transform"
	"wms-proxy/pkg/wms"
)

func TestTranslateWMSToArcGIS(t *testing.T) {
	wmsParams := &wms.WMSParams{
		Service:     "WMS",
		Version:     "1.1.1",
		Request:     "GetMap",
		Layers:      "17",
		Format:      "image/png",
		Transparent: "true",
		SRS:         "EPSG:3857",
		BBOX:        "-8238310.24,4969803.4,-8238016.75,4970096.9",
		Width:       256,
		Height:      256,
	}

	arcgisParams, err := TranslateWMSToArcGIS(wmsParams)
	if err != nil {
		t.Fatalf("TranslateWMSToArcGIS failed: %v", err)
	}

	// Verify basic translation
	if arcgisParams.BBOX != wmsParams.BBOX {
		t.Errorf("BBOX not preserved: got %s, expected %s", arcgisParams.BBOX, wmsParams.BBOX)
	}

	if arcgisParams.Size != "256,256" {
		t.Errorf("Size incorrect: got %s, expected 256,256", arcgisParams.Size)
	}

	if arcgisParams.Format != "png32" {
		t.Errorf("Format incorrect: got %s, expected png32", arcgisParams.Format)
	}

	if arcgisParams.BBoxSR != "EPSG:3857" {
		t.Errorf("BBoxSR incorrect: got %s, expected EPSG:3857", arcgisParams.BBoxSR)
	}

	if arcgisParams.Layers != "show:17" {
		t.Errorf("Layers incorrect: got %s, expected show:17", arcgisParams.Layers)
	}
}

func TestTranslateWMSToArcGISWithTransform(t *testing.T) {
	transformer := transform.NewCoordinateTransformer()

	tests := []struct {
		name            string
		wmsParams       *wms.WMSParams
		expectTransform bool
	}{
		{
			name: "Transform EPSG:3857 to EPSG:3424 - client sends EPSG:3857, backend expects EPSG:3424",
			wmsParams: &wms.WMSParams{
				Service:     "WMS",
				Version:     "1.1.1",
				Request:     "GetMap",
				Layers:      "17",
				Format:      "image/png",
				Transparent: "true",
				SRS:         "EPSG:3857",                                   // Client sends coordinates in EPSG:3857
				BBOX:        "-8238310.24,4969803.4,-8238016.75,4970096.9", // These are Web Mercator coordinates
				Width:       256,
				Height:      256,
			},
			expectTransform: true,
		},
		{
			name: "No transformation needed - client sends EPSG:3424, backend expects EPSG:3424",
			wmsParams: &wms.WMSParams{
				Service:     "WMS",
				Version:     "1.1.1",
				Request:     "GetMap",
				Layers:      "17",
				Format:      "image/png",
				Transparent: "true",
				SRS:         "EPSG:3424",                   // Client sends coordinates in EPSG:3424
				BBOX:        "629066,684288,629793,685020", // These are New Jersey State Plane coordinates
				Width:       256,
				Height:      256,
			},
			expectTransform: false,
		},
		{
			name: "Transform to WGS84",
			wmsParams: &wms.WMSParams{
				Service:     "WMS",
				Version:     "1.1.1",
				Request:     "GetMap",
				Layers:      "17",
				Format:      "image/png",
				Transparent: "true",
				SRS:         "EPSG:4326",                                   // This is the target SRS that the client wants
				BBOX:        "-8238310.24,4969803.4,-8238016.75,4970096.9", // These are Web Mercator coordinates
				Width:       256,
				Height:      256,
			},
			expectTransform: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Test with transformer
			arcgisParamsWithTransform, err := TranslateWMSToArcGISWithTransform(test.wmsParams, transformer)
			if err != nil {
				t.Fatalf("TranslateWMSToArcGISWithTransform failed: %v", err)
			}

			// Test without transformer for comparison
			arcgisParamsWithoutTransform, err := TranslateWMSToArcGIS(test.wmsParams)
			if err != nil {
				t.Fatalf("TranslateWMSToArcGIS failed: %v", err)
			}

			if test.expectTransform {
				// BBOX should be different when transformation is expected
				if arcgisParamsWithTransform.BBOX == test.wmsParams.BBOX {
					t.Error("Expected BBOX to be transformed, but it remained the same")
				}
				t.Logf("Original BBOX: %s", test.wmsParams.BBOX)
				t.Logf("Transformed BBOX: %s", arcgisParamsWithTransform.BBOX)
			} else {
				// BBOX should be the same when no transformation is expected
				if arcgisParamsWithTransform.BBOX != test.wmsParams.BBOX {
					t.Errorf("Expected BBOX to remain the same, but it was transformed: %s -> %s",
						test.wmsParams.BBOX, arcgisParamsWithTransform.BBOX)
				}
			}

			// Other parameters should be the same regardless of transformation
			if arcgisParamsWithTransform.Size != arcgisParamsWithoutTransform.Size {
				t.Error("Size parameter differs between transformer and non-transformer versions")
			}

			if arcgisParamsWithTransform.Format != arcgisParamsWithoutTransform.Format {
				t.Error("Format parameter differs between transformer and non-transformer versions")
			}

			if arcgisParamsWithTransform.Layers != arcgisParamsWithoutTransform.Layers {
				t.Error("Layers parameter differs between transformer and non-transformer versions")
			}
		})
	}
}

func TestTranslateFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"image/png", "png32"},
		{"png", "png32"},
		{"image/jpeg", "jpg"},
		{"jpeg", "jpg"},
		{"jpg", "jpg"},
		{"image/gif", "gif"},
		{"gif", "gif"},
		{"unknown", "png32"}, // Default
	}

	for _, test := range tests {
		result := translateFormat(test.input)
		if result != test.expected {
			t.Errorf("translateFormat(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestTranslateSRS(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"EPSG:3857", "EPSG:3857"},
		{"3857", "EPSG:3857"},
		{"900913", "EPSG:3857"},
		{"EPSG:4326", "EPSG:4326"},
		{"4326", "EPSG:4326"},
		{"EPSG:3424", "EPSG:3424"},
		{"3424", "EPSG:3424"},
		{"unknown", "unknown"},
		{"", "EPSG:3857"}, // Default
	}

	for _, test := range tests {
		result := translateSRS(test.input)
		if result != test.expected {
			t.Errorf("translateSRS(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestTranslateLayers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"17", "show:17"},
		{"1,2,3", "show:1,show:2,show:3"},
		{"", "show:0"},      // Default
		{" 17 ", "show:17"}, // Trimmed
		{"17,", "show:17"},  // Trailing comma
	}

	for _, test := range tests {
		result := translateLayers(test.input)
		if result != test.expected {
			t.Errorf("translateLayers(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestTranslateTransparent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true", "true"},
		{"TRUE", "true"},
		{"1", "true"},
		{"yes", "true"},
		{"false", "false"},
		{"FALSE", "false"},
		{"0", "false"},
		{"no", "false"},
		{"", "true"},        // Default
		{"unknown", "true"}, // Default
	}

	for _, test := range tests {
		result := translateTransparent(test.input)
		if result != test.expected {
			t.Errorf("translateTransparent(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestBuildArcGISURL(t *testing.T) {
	baseURL := "https://example.com"
	servicePath := "/arcgis/rest/services/test/MapServer/export"
	params := &wms.ArcGISParams{
		BBOX:        "-8238310.24,4969803.4,-8238016.75,4970096.9",
		Size:        "256,256",
		Format:      "png32",
		BBoxSR:      "EPSG:3857",
		ImageSR:     "EPSG:3857",
		Layers:      "show:17",
		Transparent: "true",
		DPI:         96,
		F:           "image",
	}

	result := BuildArcGISURL(baseURL, servicePath, params)

	expectedBase := baseURL + servicePath
	if !strings.HasPrefix(result, expectedBase) {
		t.Errorf("URL doesn't start with expected base: %s", result)
	}

	// Check that all parameters are present
	expectedParams := []string{
		"bbox=-8238310.24%2C4969803.4%2C-8238016.75%2C4970096.9",
		"size=256%2C256",
		"format=png32",
		"bboxSR=EPSG%3A3857",
		"imageSR=EPSG%3A3857",
		"layers=show%3A17",
		"transparent=true",
		"dpi=96",
		"f=image",
	}

	for _, param := range expectedParams {
		if !strings.Contains(result, param) {
			t.Errorf("URL missing parameter: %s\nFull URL: %s", param, result)
		}
	}

	t.Logf("Generated URL: %s", result)
}
