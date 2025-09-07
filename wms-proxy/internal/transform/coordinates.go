package transform

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// CoordinateTransformer handles coordinate transformations between different spatial reference systems
type CoordinateTransformer struct {
	transformers map[string]map[string]TransformFunc
}

// TransformFunc represents a function that transforms coordinates from one CRS to another
type TransformFunc func(x, y float64) (float64, float64, error)

// BBox represents a bounding box with min/max x/y coordinates
type BBox struct {
	MinX, MinY, MaxX, MaxY float64
}

// NewCoordinateTransformer creates a new coordinate transformer with predefined transformations
func NewCoordinateTransformer() *CoordinateTransformer {
	ct := &CoordinateTransformer{
		transformers: make(map[string]map[string]TransformFunc),
	}

	// Initialize predefined transformations
	ct.initializeTransformations()

	return ct
}

// initializeTransformations sets up the supported coordinate transformations
func (ct *CoordinateTransformer) initializeTransformations() {
	// EPSG:3857 (Web Mercator) to EPSG:3424 (NAD83 / New Jersey (ftUS))
	ct.addTransformation("EPSG:3857", "EPSG:3424", webMercatorToNAD83NewJersey)

	// EPSG:3857 to EPSG:4326 (WGS84 Geographic)
	ct.addTransformation("EPSG:3857", "EPSG:4326", webMercatorToWGS84)

	// EPSG:4326 (WGS84 Geographic) to EPSG:3424 (NAD83 / New Jersey (ftUS))
	ct.addTransformation("EPSG:4326", "EPSG:3424", wgs84ToNAD83NewJersey)

	// Reverse transformations
	// EPSG:3424 to EPSG:3857 (reverse of webMercatorToNAD83NewJersey)
	ct.addTransformation("EPSG:3424", "EPSG:3857", nad83NewJerseyToWebMercator)

	// EPSG:4326 to EPSG:3857 (reverse of webMercatorToWGS84)
	ct.addTransformation("EPSG:4326", "EPSG:3857", wgs84ToWebMercator)

	// Add identity transformations (same CRS)
	ct.addTransformation("EPSG:3857", "EPSG:3857", identityTransform)
	ct.addTransformation("EPSG:3424", "EPSG:3424", identityTransform)
	ct.addTransformation("EPSG:4326", "EPSG:4326", identityTransform)
}

// addTransformation adds a transformation function between two coordinate systems
func (ct *CoordinateTransformer) addTransformation(fromCRS, toCRS string, transformFunc TransformFunc) {
	if ct.transformers[fromCRS] == nil {
		ct.transformers[fromCRS] = make(map[string]TransformFunc)
	}
	ct.transformers[fromCRS][toCRS] = transformFunc
}

// TransformBBox transforms a bounding box from one CRS to another
func (ct *CoordinateTransformer) TransformBBox(bboxStr, fromCRS, toCRS string) (string, error) {
	// Parse the bbox string
	bbox, err := parseBBox(bboxStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse bbox: %w", err)
	}

	// Get the transformation function
	transformFunc, err := ct.getTransformFunc(fromCRS, toCRS)
	if err != nil {
		return "", err
	}

	// Transform the corner coordinates
	minX, minY, err := transformFunc(bbox.MinX, bbox.MinY)
	if err != nil {
		return "", fmt.Errorf("failed to transform min coordinates: %w", err)
	}

	maxX, maxY, err := transformFunc(bbox.MaxX, bbox.MaxY)
	if err != nil {
		return "", fmt.Errorf("failed to transform max coordinates: %w", err)
	}

	// Return the transformed bbox as a string
	return fmt.Sprintf("%.6f,%.6f,%.6f,%.6f", minX, minY, maxX, maxY), nil
}

// getTransformFunc retrieves the transformation function for the given CRS pair
func (ct *CoordinateTransformer) getTransformFunc(fromCRS, toCRS string) (TransformFunc, error) {
	// Normalize CRS names
	fromCRS = normalizeCRS(fromCRS)
	toCRS = normalizeCRS(toCRS)

	if ct.transformers[fromCRS] == nil {
		return nil, fmt.Errorf("unsupported source CRS: %s", fromCRS)
	}

	transformFunc, exists := ct.transformers[fromCRS][toCRS]
	if !exists {
		return nil, fmt.Errorf("unsupported transformation from %s to %s", fromCRS, toCRS)
	}

	return transformFunc, nil
}

// NormalizeCRS normalizes CRS names to a standard format (public method)
func (ct *CoordinateTransformer) NormalizeCRS(crs string) string {
	return normalizeCRS(crs)
}

// normalizeCRS normalizes CRS names to a standard format
func normalizeCRS(crs string) string {
	originalCRS := strings.TrimSpace(crs)
	upperCRS := strings.ToUpper(originalCRS)

	// Handle common variations
	switch upperCRS {
	case "3857", "900913":
		return "EPSG:3857"
	case "4326":
		return "EPSG:4326"
	case "3424", "102711":
		return "EPSG:3424"
	case "EPSG:3857", "EPSG:900913":
		return "EPSG:3857"
	case "EPSG:4326":
		return "EPSG:4326"
	case "EPSG:3424", "EPSG:102711":
		return "EPSG:3424"
	default:
		// If it's just a number, add EPSG: prefix
		if _, err := strconv.Atoi(originalCRS); err == nil {
			return "EPSG:" + originalCRS
		}
		// Return original case for unrecognized values
		return originalCRS
	}
}

// parseBBox parses a bbox string in the format "minx,miny,maxx,maxy"
func parseBBox(bboxStr string) (*BBox, error) {
	parts := strings.Split(bboxStr, ",")
	if len(parts) != 4 {
		return nil, fmt.Errorf("bbox must have 4 comma-separated values, got %d", len(parts))
	}

	coords := make([]float64, 4)
	for i, part := range parts {
		coord, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid coordinate at position %d: %s", i, part)
		}
		coords[i] = coord
	}

	return &BBox{
		MinX: coords[0],
		MinY: coords[1],
		MaxX: coords[2],
		MaxY: coords[3],
	}, nil
}

// identityTransform returns the same coordinates (no transformation)
func identityTransform(x, y float64) (float64, float64, error) {
	return x, y, nil
}

// webMercatorToWGS84 converts Web Mercator (EPSG:3857) to WGS84 Geographic (EPSG:4326)
func webMercatorToWGS84(x, y float64) (float64, float64, error) {
	const earthRadius = 6378137.0 // Earth radius in meters

	// Convert to longitude/latitude
	lon := x / earthRadius * 180.0 / math.Pi
	lat := math.Atan(math.Exp(y/earthRadius))*360.0/math.Pi - 90.0

	return lon, lat, nil
}

// webMercatorToNAD83NewJersey converts Web Mercator (EPSG:3857) to NAD83 / New Jersey (ftUS) (EPSG:3424)
// This is a simplified implementation that goes through WGS84 as an intermediate step
func webMercatorToNAD83NewJersey(x, y float64) (float64, float64, error) {
	// First convert Web Mercator to WGS84
	lon, lat, err := webMercatorToWGS84(x, y)
	if err != nil {
		return 0, 0, err
	}

	// Then convert WGS84 to NAD83 New Jersey
	// This is a simplified transformation - in a production system, you would use
	// a proper geodetic transformation library
	return wgs84ToNAD83NewJersey(lon, lat)
}

// wgs84ToNAD83NewJersey converts WGS84 Geographic to NAD83 / New Jersey (ftUS)
// EPSG:3424 parameters: +proj=tmerc +lat_0=38.8333333333333 +lon_0=-74.5 +k=0.9999 +x_0=492125 +y_0=0 +ellps=GRS80 +units=us-ft +no_defs
func wgs84ToNAD83NewJersey(lon, lat float64) (float64, float64, error) {
	// Transverse Mercator projection parameters for EPSG:3424
	const (
		lat0           = 38.8333333333333 * math.Pi / 180.0 // Latitude of origin in radians
		lon0           = -74.5 * math.Pi / 180.0            // Central meridian in radians
		k0             = 0.9999                             // Scale factor
		falseE         = 492125.0                           // False easting in US survey feet
		falseN         = 0.0                                // False northing in US survey feet
		a              = 6378137.0                          // Semi-major axis (GRS80)
		f              = 1.0 / 298.257222101                // Flattening (GRS80)
		metersToUSFeet = 3.280833333                        // Conversion factor from meters to US survey feet
	)

	// Convert degrees to radians
	latRad := lat * math.Pi / 180.0
	lonRad := lon * math.Pi / 180.0

	// Calculate eccentricity
	e := math.Sqrt(2*f - f*f)
	e2 := e * e

	// Calculate delta longitude
	deltaLon := lonRad - lon0

	// Simplified Transverse Mercator formulas (suitable for small areas like New Jersey)
	// For production use, implement the full Transverse Mercator projection

	// Calculate N (radius of curvature in the prime vertical)
	sinLat := math.Sin(latRad)
	N := a / math.Sqrt(1-e2*sinLat*sinLat)

	// Calculate T, C, A
	tanLat := math.Tan(latRad)
	T := tanLat * tanLat
	cosLat := math.Cos(latRad)
	C := e2 * cosLat * cosLat / (1 - e2)
	A := cosLat * deltaLon

	// Calculate M (meridional arc)
	e4 := e2 * e2
	e6 := e4 * e2
	M0 := a * ((1-e2/4-3*e4/64-5*e6/256)*lat0 -
		(3*e2/8+3*e4/32+45*e6/1024)*math.Sin(2*lat0) +
		(15*e4/256+45*e6/1024)*math.Sin(4*lat0) -
		(35*e6/3072)*math.Sin(6*lat0))

	M := a * ((1-e2/4-3*e4/64-5*e6/256)*latRad -
		(3*e2/8+3*e4/32+45*e6/1024)*math.Sin(2*latRad) +
		(15*e4/256+45*e6/1024)*math.Sin(4*latRad) -
		(35*e6/3072)*math.Sin(6*latRad))

	// Calculate easting and northing in meters
	A2 := A * A
	A4 := A2 * A2
	A6 := A4 * A2

	x := k0 * N * (A + (1-T+C)*A*A2/6 + (5-18*T+T*T+72*C-58*e2)*A*A4/120)
	y := k0 * (M - M0 + N*tanLat*(A2/2+(5-T+9*C+4*C*C)*A4/24+(61-58*T+T*T+600*C-330*e2)*A6/720))

	// Convert to US survey feet and apply false easting/northing
	eastingFt := x*metersToUSFeet + falseE
	northingFt := y*metersToUSFeet + falseN

	return eastingFt, northingFt, nil
}

// nad83NewJerseyToWebMercator converts NAD83 / New Jersey (ftUS) (EPSG:3424) to Web Mercator (EPSG:3857)
// This is the reverse of webMercatorToNAD83NewJersey
func nad83NewJerseyToWebMercator(eastingFt, northingFt float64) (float64, float64, error) {
	// First convert from NAD83 New Jersey to WGS84
	lon, lat, err := nad83NewJerseyToWGS84(eastingFt, northingFt)
	if err != nil {
		return 0, 0, err
	}

	// Then convert WGS84 to Web Mercator
	return wgs84ToWebMercator(lon, lat)
}

// nad83NewJerseyToWGS84 converts NAD83 / New Jersey (ftUS) to WGS84 Geographic
// This is the reverse of wgs84ToNAD83NewJersey
func nad83NewJerseyToWGS84(eastingFt, northingFt float64) (float64, float64, error) {
	// Transverse Mercator projection parameters for EPSG:3424
	const (
		lat0           = 38.8333333333333 * math.Pi / 180.0 // Latitude of origin in radians
		lon0           = -74.5 * math.Pi / 180.0            // Central meridian in radians
		k0             = 0.9999                             // Scale factor
		falseE         = 492125.0                           // False easting in US survey feet
		falseN         = 0.0                                // False northing in US survey feet
		a              = 6378137.0                          // Semi-major axis (GRS80)
		f              = 1.0 / 298.257222101                // Flattening (GRS80)
		metersToUSFeet = 3.280833333                        // Conversion factor from meters to US survey feet
	)

	// Convert from US survey feet to meters and remove false easting/northing
	x := (eastingFt - falseE) / metersToUSFeet
	y := (northingFt - falseN) / metersToUSFeet

	// Calculate eccentricity
	e := math.Sqrt(2*f - f*f)
	e2 := e * e

	// This is a simplified reverse transformation
	// For production use, implement the full inverse Transverse Mercator projection

	// Calculate M0 (meridional arc at origin)
	e4 := e2 * e2
	e6 := e4 * e2
	M0 := a * ((1-e2/4-3*e4/64-5*e6/256)*lat0 -
		(3*e2/8+3*e4/32+45*e6/1024)*math.Sin(2*lat0) +
		(15*e4/256+45*e6/1024)*math.Sin(4*lat0) -
		(35*e6/3072)*math.Sin(6*lat0))

	// Calculate latitude (simplified)
	M := M0 + y/k0
	mu := M / (a * (1 - e2/4 - 3*e4/64 - 5*e6/256))

	// Footprint latitude (simplified calculation)
	lat1 := mu + (3*e2/2-27*e4/32)*math.Sin(2*mu) +
		(21*e4/16-55*e6/32)*math.Sin(4*mu) +
		(151*e6/96)*math.Sin(6*mu)

	// Calculate longitude (simplified)
	N1 := a / math.Sqrt(1-e2*math.Sin(lat1)*math.Sin(lat1))
	T1 := math.Tan(lat1) * math.Tan(lat1)
	C1 := e2 * math.Cos(lat1) * math.Cos(lat1) / (1 - e2)
	R1 := a * (1 - e2) / math.Pow(1-e2*math.Sin(lat1)*math.Sin(lat1), 1.5)
	D := x / (N1 * k0)

	lat := lat1 - (N1*math.Tan(lat1)/R1)*(D*D/2-(5+3*T1+10*C1-4*C1*C1-9*e2)*D*D*D*D/24)
	lon := lon0 + (D-(1+2*T1+C1)*D*D*D/6+(5-2*C1+28*T1-3*C1*C1+8*e2+24*T1*T1)*D*D*D*D*D/120)/math.Cos(lat1)

	// Convert from radians to degrees
	latDeg := lat * 180.0 / math.Pi
	lonDeg := lon * 180.0 / math.Pi

	return lonDeg, latDeg, nil
}

// wgs84ToWebMercator converts WGS84 Geographic (EPSG:4326) to Web Mercator (EPSG:3857)
// This is the reverse of webMercatorToWGS84
func wgs84ToWebMercator(lon, lat float64) (float64, float64, error) {
	const earthRadius = 6378137.0 // Earth radius in meters

	// Convert to Web Mercator
	x := lon * math.Pi / 180.0 * earthRadius
	y := math.Log(math.Tan((90.0+lat)*math.Pi/360.0)) * earthRadius

	return x, y, nil
}
