package main

import (
	"fmt"
	"log"
	"os"

	"wms-proxy/internal/config"
	"wms-proxy/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Print startup information
	fmt.Printf("WMS Proxy Server starting...\n")
	fmt.Printf("Proxy Port: %d\n", cfg.ProxyPort)
	fmt.Printf("HTTPS Enabled: %t\n", cfg.EnableHTTPS)
	if cfg.EnableHTTPS {
		fmt.Printf("Certificate File: %s\n", cfg.CertFile)
		fmt.Printf("Key File: %s\n", cfg.KeyFile)
	}
	fmt.Printf("ArcGIS Host: %s\n", cfg.ArcGISHost)
	fmt.Printf("ArcGIS Scheme: %s\n", cfg.ArcGISScheme)
	fmt.Printf("ArcGIS Service: %s\n", cfg.ArcGISService)
	fmt.Printf("Request Timeout: %s\n", cfg.RequestTimeout)
	fmt.Printf("Log Level: %s\n", cfg.LogLevel)

	// Create and start server
	srv := server.New(cfg)

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
