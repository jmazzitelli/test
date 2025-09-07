package services

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"wms-proxy/internal/client"
)

// BackendSRDetector manages detection of backend spatial reference systems
type BackendSRDetector struct {
	arcgisClient client.ArcGISClientInterface
	logger       *slog.Logger
	cache        map[string]string // servicePath -> EPSG code
	cacheMutex   sync.RWMutex
	cacheTTL     time.Duration
	cacheExpiry  map[string]time.Time
}

// NewBackendSRDetector creates a new backend spatial reference detector
func NewBackendSRDetector(arcgisClient client.ArcGISClientInterface, logger *slog.Logger) *BackendSRDetector {
	return &BackendSRDetector{
		arcgisClient: arcgisClient,
		logger:       logger,
		cache:        make(map[string]string),
		cacheExpiry:  make(map[string]time.Time),
		cacheTTL:     15 * time.Minute, // Cache for 15 minutes
	}
}

// GetBackendSR detects the spatial reference system expected by the backend service
func (d *BackendSRDetector) GetBackendSR(ctx context.Context, servicePath string) (string, error) {
	// Check cache first
	d.cacheMutex.RLock()
	if cachedSR, exists := d.cache[servicePath]; exists {
		if expiry, hasExpiry := d.cacheExpiry[servicePath]; hasExpiry && time.Now().Before(expiry) {
			d.cacheMutex.RUnlock()
			d.logger.Debug("Using cached backend SR", "service_path", servicePath, "sr", cachedSR)
			return cachedSR, nil
		}
	}
	d.cacheMutex.RUnlock()

	// Query service metadata
	d.logger.Info("Querying backend service metadata", "service_path", servicePath)
	metadata, err := d.arcgisClient.GetServiceMetadata(ctx, servicePath)
	if err != nil {
		return "", fmt.Errorf("failed to get service metadata: %w", err)
	}

	// Extract spatial reference - prefer LatestWKID when available as it represents the modern standard
	var backendSR string
	if metadata.SpatialReference.LatestWKID != 0 {
		backendSR = "EPSG:" + strconv.Itoa(metadata.SpatialReference.LatestWKID)
	} else if metadata.SpatialReference.WKID != 0 {
		backendSR = "EPSG:" + strconv.Itoa(metadata.SpatialReference.WKID)
	} else {
		// Fallback to a reasonable default for New Jersey services
		backendSR = "EPSG:3424"
		d.logger.Warn("Could not determine backend SR from metadata, using fallback",
			"service_path", servicePath,
			"fallback_sr", backendSR,
			"metadata_wkid", metadata.SpatialReference.WKID,
			"metadata_latest_wkid", metadata.SpatialReference.LatestWKID)
	}

	// Cache the result
	d.cacheMutex.Lock()
	d.cache[servicePath] = backendSR
	d.cacheExpiry[servicePath] = time.Now().Add(d.cacheTTL)
	d.cacheMutex.Unlock()

	d.logger.Info("Detected backend spatial reference system",
		"service_path", servicePath,
		"backend_sr", backendSR,
		"wkid", metadata.SpatialReference.WKID,
		"latest_wkid", metadata.SpatialReference.LatestWKID)

	return backendSR, nil
}

// ClearCache clears the spatial reference cache
func (d *BackendSRDetector) ClearCache() {
	d.cacheMutex.Lock()
	defer d.cacheMutex.Unlock()

	d.cache = make(map[string]string)
	d.cacheExpiry = make(map[string]time.Time)
	d.logger.Info("Backend SR cache cleared")
}

// GetCacheStats returns cache statistics for monitoring
func (d *BackendSRDetector) GetCacheStats() map[string]interface{} {
	d.cacheMutex.RLock()
	defer d.cacheMutex.RUnlock()

	validEntries := 0
	expiredEntries := 0
	now := time.Now()

	for servicePath, expiry := range d.cacheExpiry {
		if _, exists := d.cache[servicePath]; exists {
			if now.Before(expiry) {
				validEntries++
			} else {
				expiredEntries++
			}
		}
	}

	return map[string]interface{}{
		"total_entries":     len(d.cache),
		"valid_entries":     validEntries,
		"expired_entries":   expiredEntries,
		"cache_ttl_minutes": d.cacheTTL.Minutes(),
	}
}
