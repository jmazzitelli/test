package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the proxy server
type Config struct {
	ArcGISHost     string
	ArcGISScheme   string
	ArcGISService  string
	ProxyPort      int
	RequestTimeout time.Duration
	LogLevel       string
	EnableHTTPS    bool
	CertFile       string
	KeyFile        string
}

// Load reads configuration from environment variables with sensible defaults
func Load() (*Config, error) {
	cfg := &Config{
		ArcGISHost:     getEnvString("ARCGIS_HOST", "localhost"),
		ArcGISScheme:   getEnvString("ARCGIS_SCHEME", "https"),
		ArcGISService:  getEnvString("ARCGIS_SERVICE", "/arcgis/rest/services/Features/Environmental_admin/MapServer/export"),
		ProxyPort:      getEnvInt("PROXY_PORT", 8080),
		RequestTimeout: time.Duration(getEnvInt("REQUEST_TIMEOUT", 30)) * time.Second,
		LogLevel:       getEnvString("LOG_LEVEL", "info"),
		EnableHTTPS:    getEnvBool("ENABLE_HTTPS", false),
		CertFile:       getEnvString("CERT_FILE", "/app/certs/server.crt"),
		KeyFile:        getEnvString("KEY_FILE", "/app/certs/server.key"),
	}

	// Validate required configuration
	if cfg.ArcGISHost == "" {
		return nil, fmt.Errorf("ARCGIS_HOST is required")
	}

	if cfg.ArcGISScheme != "http" && cfg.ArcGISScheme != "https" {
		return nil, fmt.Errorf("ARCGIS_SCHEME must be 'http' or 'https'")
	}

	if cfg.ProxyPort <= 0 || cfg.ProxyPort > 65535 {
		return nil, fmt.Errorf("PROXY_PORT must be between 1 and 65535")
	}

	// Validate HTTPS configuration
	if cfg.EnableHTTPS {
		if cfg.CertFile == "" {
			return nil, fmt.Errorf("CERT_FILE is required when HTTPS is enabled")
		}
		if cfg.KeyFile == "" {
			return nil, fmt.Errorf("KEY_FILE is required when HTTPS is enabled")
		}
	}

	return cfg, nil
}

// GetArcGISBaseURL returns the base URL for the ArcGIS server
func (c *Config) GetArcGISBaseURL() string {
	return fmt.Sprintf("%s://%s", c.ArcGISScheme, c.ArcGISHost)
}

// GetProxyAddress returns the address the proxy should listen on
func (c *Config) GetProxyAddress() string {
	return fmt.Sprintf(":%d", c.ProxyPort)
}

func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	return defaultValue
}
