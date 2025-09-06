package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"wms-proxy/internal/client"
	"wms-proxy/internal/config"
	"wms-proxy/internal/handlers"
)

// Server represents the WMS proxy server
type Server struct {
	config       *config.Config
	logger       *slog.Logger
	httpServer   *http.Server
	arcgisClient *client.ArcGISClient
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	// Setup logger
	logger := setupLogger(cfg.LogLevel)

	// Create ArcGIS client
	arcgisClient := client.NewArcGISClient(cfg.GetArcGISBaseURL(), cfg.RequestTimeout)

	return &Server{
		config:       cfg,
		logger:       logger,
		arcgisClient: arcgisClient,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Setup routes
	router := s.setupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         s.config.GetProxyAddress(),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		protocol := "HTTP"
		if s.config.EnableHTTPS {
			protocol = "HTTPS"
		}

		s.logger.Info("Starting WMS proxy server",
			"protocol", protocol,
			"address", s.config.GetProxyAddress(),
			"arcgis_host", s.config.ArcGISHost,
			"https_enabled", s.config.EnableHTTPS,
		)

		var err error
		if s.config.EnableHTTPS {
			err = s.httpServer.ListenAndServeTLS(s.config.CertFile, s.config.KeyFile)
		} else {
			err = s.httpServer.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	return s.waitForShutdown()
}

// setupRoutes configures the HTTP routes
func (s *Server) setupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Add logging middleware
	router.Use(s.loggingMiddleware)

	// Health check endpoint
	healthHandler := handlers.NewHealthHandler(s.arcgisClient)
	router.Handle("/health", healthHandler).Methods("GET")

	// ArcGIS REST API proxy (direct passthrough)
	arcgisProxyHandler := handlers.NewArcGISProxyHandler(s.arcgisClient, s.logger, s.config.GetArcGISBaseURL())

	// Handle ArcGIS REST API paths directly
	router.PathPrefix("/arcgis/").Handler(arcgisProxyHandler).Methods("GET")

	// WMS endpoints (for WMS clients)
	wmsHandler := handlers.NewWMSHandler(s.arcgisClient, s.logger, s.config.GetArcGISBaseURL(), s.config.ArcGISService)

	// Handle WMS requests
	router.Handle("/wms", wmsHandler).Methods("GET")

	// Root path defaults to WMS for backward compatibility
	router.Handle("/", wmsHandler).Methods("GET")

	return router
}

// loggingMiddleware logs HTTP requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		s.logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// waitForShutdown waits for interrupt signal and gracefully shuts down the server
func (s *Server) waitForShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	s.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	s.logger.Info("Server exited")
	return nil
}

// setupLogger creates a structured logger based on log level
func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level

	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}
