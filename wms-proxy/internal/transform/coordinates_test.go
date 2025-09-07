package transform

import (
	"math"
	"testing"
)

func TestNewCoordinateTransformer(t *testing.T) {
	transformer := NewCoordinateTransformer()
	if transformer == nil {
		t.Fatal("NewCoordinateTransformer returned nil")
	}

	if transformer.transformers == nil {
		t.Fatal("transformers map is nil")
	}

	// Check that some basic transformations are initialized
	if transformer.transformers["EPSG:3857"] == nil {
		t.Error("EPSG:3857 transformations not initialized")
	}
}

func TestNormalizeCRS(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"3857", "EPSG:3857"},
		{"4326", "EPSG:4326"},
		{"3424", "EPSG:3424"},
		{"102711", "EPSG:3424"}, // ESRI code that maps to EPSG:3424
		{"EPSG:3857", "EPSG:3857"},
		{"epsg:4326", "EPSG:4326"},
		{"EPSG:102711", "EPSG:3424"}, // ESRI code that maps to EPSG:3424
		{"900913", "EPSG:3857"},
		{"", ""},
		{"invalid", "invalid"},
	}

	for _, test := range tests {
		result := normalizeCRS(test.input)
		if result != test.expected {
			t.Errorf("normalizeCRS(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestParseBBox(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
		expected    *BBox
	}{
		{
			input:       "-8238310.24,4969803.4,-8238016.75,4970096.9",
			expectError: false,
			expected: &BBox{
				MinX: -8238310.24,
				MinY: 4969803.4,
				MaxX: -8238016.75,
				MaxY: 4970096.9,
			},
		},
		{
			input:       "0,0,100,100",
			expectError: false,
			expected: &BBox{
				MinX: 0,
				MinY: 0,
				MaxX: 100,
				MaxY: 100,
			},
		},
		{
			input:       "1,2,3", // Only 3 values
			expectError: true,
			expected:    nil,
		},
		{
			input:       "1,2,3,4,5", // Too many values
			expectError: true,
			expected:    nil,
		},
		{
			input:       "1,invalid,3,4", // Invalid number
			expectError: true,
			expected:    nil,
		},
		{
			input:       "", // Empty string
			expectError: true,
			expected:    nil,
		},
	}

	for _, test := range tests {
		result, err := parseBBox(test.input)

		if test.expectError {
			if err == nil {
				t.Errorf("parseBBox(%q) expected error but got none", test.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("parseBBox(%q) unexpected error: %v", test.input, err)
			continue
		}

		if result.MinX != test.expected.MinX ||
			result.MinY != test.expected.MinY ||
			result.MaxX != test.expected.MaxX ||
			result.MaxY != test.expected.MaxY {
			t.Errorf("parseBBox(%q) = %+v, expected %+v", test.input, result, test.expected)
		}
	}
}

func TestIdentityTransform(t *testing.T) {
	x, y := 123.456, 789.012
	resultX, resultY, err := identityTransform(x, y)

	if err != nil {
		t.Errorf("identityTransform unexpected error: %v", err)
	}

	if resultX != x || resultY != y {
		t.Errorf("identityTransform(%f, %f) = (%f, %f), expected (%f, %f)",
			x, y, resultX, resultY, x, y)
	}
}

func TestWebMercatorToWGS84(t *testing.T) {
	tests := []struct {
		name                     string
		x, y                     float64
		expectedLon, expectedLat float64
		tolerance                float64
	}{
		{
			name: "New York area",
			x:    -8238310.24, y: 4969803.4,
			expectedLon: -74.006, expectedLat: 40.711,
			tolerance: 0.01, // Allow 0.01 degree tolerance
		},
		{
			name: "Origin",
			x:    0, y: 0,
			expectedLon: 0, expectedLat: 0,
			tolerance: 0.001,
		},
		{
			name: "London area",
			x:    -14000, y: 6711000,
			expectedLon: -0.126, expectedLat: 51.5,
			tolerance: 0.1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lon, lat, err := webMercatorToWGS84(test.x, test.y)
			if err != nil {
				t.Errorf("webMercatorToWGS84 unexpected error: %v", err)
				return
			}

			if math.Abs(lon-test.expectedLon) > test.tolerance {
				t.Errorf("longitude: got %f, expected %f (tolerance %f)",
					lon, test.expectedLon, test.tolerance)
			}

			if math.Abs(lat-test.expectedLat) > test.tolerance {
				t.Errorf("latitude: got %f, expected %f (tolerance %f)",
					lat, test.expectedLat, test.tolerance)
			}
		})
	}
}

func TestWebMercatorToNAD83NewJersey(t *testing.T) {
	// Test with New Jersey coordinates
	x, y := -8238310.24, 4969803.4 // Web Mercator coordinates for New Jersey area

	eastingFt, northingFt, err := webMercatorToNAD83NewJersey(x, y)
	if err != nil {
		t.Errorf("webMercatorToNAD83NewJersey unexpected error: %v", err)
		return
	}

	// Check that we get reasonable values for New Jersey State Plane coordinates
	// New Jersey State Plane coordinates should be positive and in a reasonable range
	if eastingFt < 0 || eastingFt > 1000000 {
		t.Errorf("easting %f is outside reasonable range for New Jersey", eastingFt)
	}

	if northingFt < 0 || northingFt > 1000000 {
		t.Errorf("northing %f is outside reasonable range for New Jersey", northingFt)
	}

	t.Logf("Transformed coordinates: easting=%f ft, northing=%f ft", eastingFt, northingFt)
}

func TestTransformBBox(t *testing.T) {
	transformer := NewCoordinateTransformer()

	tests := []struct {
		name        string
		bbox        string
		fromCRS     string
		toCRS       string
		expectError bool
	}{
		{
			name:        "Web Mercator to NAD83 New Jersey",
			bbox:        "-8238310.24,4969803.4,-8238016.75,4970096.9",
			fromCRS:     "EPSG:3857",
			toCRS:       "EPSG:3424",
			expectError: false,
		},
		{
			name:        "Web Mercator to WGS84",
			bbox:        "-8238310.24,4969803.4,-8238016.75,4970096.9",
			fromCRS:     "EPSG:3857",
			toCRS:       "EPSG:4326",
			expectError: false,
		},
		{
			name:        "Identity transformation",
			bbox:        "-8238310.24,4969803.4,-8238016.75,4970096.9",
			fromCRS:     "EPSG:3857",
			toCRS:       "EPSG:3857",
			expectError: false,
		},
		{
			name:        "Invalid bbox format",
			bbox:        "invalid,bbox",
			fromCRS:     "EPSG:3857",
			toCRS:       "EPSG:3424",
			expectError: true,
		},
		{
			name:        "Unsupported transformation",
			bbox:        "-8238310.24,4969803.4,-8238016.75,4970096.9",
			fromCRS:     "EPSG:3857",
			toCRS:       "EPSG:9999", // Unsupported CRS
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := transformer.TransformBBox(test.bbox, test.fromCRS, test.toCRS)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == "" {
				t.Error("result is empty")
				return
			}

			// Verify the result can be parsed as a valid bbox
			_, parseErr := parseBBox(result)
			if parseErr != nil {
				t.Errorf("result bbox %q is not valid: %v", result, parseErr)
			}

			t.Logf("Transformation result: %s -> %s", test.bbox, result)
		})
	}
}

func TestTransformBBoxNormalization(t *testing.T) {
	transformer := NewCoordinateTransformer()

	// Test that CRS normalization works in TransformBBox
	tests := []struct {
		fromCRS string
		toCRS   string
	}{
		{"3857", "3424"},      // Numbers without EPSG: prefix
		{"EPSG:3857", "3424"}, // Mixed formats
		{"3857", "EPSG:3424"}, // Mixed formats
	}

	bbox := "-8238310.24,4969803.4,-8238016.75,4970096.9"

	for _, test := range tests {
		t.Run(test.fromCRS+"_to_"+test.toCRS, func(t *testing.T) {
			result, err := transformer.TransformBBox(bbox, test.fromCRS, test.toCRS)
			if err != nil {
				t.Errorf("unexpected error with CRS normalization: %v", err)
				return
			}

			if result == "" {
				t.Error("result is empty")
			}
		})
	}
}

// Benchmark tests
func BenchmarkTransformBBox(b *testing.B) {
	transformer := NewCoordinateTransformer()
	bbox := "-8238310.24,4969803.4,-8238016.75,4970096.9"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := transformer.TransformBBox(bbox, "EPSG:3857", "EPSG:3424")
		if err != nil {
			b.Fatalf("transformation failed: %v", err)
		}
	}
}

func BenchmarkWebMercatorToWGS84(b *testing.B) {
	x, y := -8238310.24, 4969803.4

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := webMercatorToWGS84(x, y)
		if err != nil {
			b.Fatalf("transformation failed: %v", err)
		}
	}
}
