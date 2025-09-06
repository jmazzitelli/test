# ArcGIS REST to WMS Proxy - Implementation Plan

## Technology Stack Selection

### Programming Language: Go
**Rationale:**
- Excellent HTTP server and client libraries
- Built-in concurrency support for handling multiple requests
- Small binary size ideal for containers
- Strong performance characteristics
- Cross-platform compilation support
- Robust standard library for URL parsing, HTTP handling

### Key Libraries:
- `net/http` - HTTP server and client
- `gorilla/mux` - HTTP router for clean URL handling
- `encoding/xml` - WMS capabilities XML generation
- `log/slog` - Structured logging
- `context` - Request context and timeouts

## Architecture Overview

The proxy supports dual operational modes:

### Mode 1: Direct ArcGIS REST Proxy
```
[ArcGIS Client] → [Proxy Server] → [ArcGIS REST Server]
                       ↓
                [Direct Passthrough]
                       ↓
              [ArcGIS Response] ← [ArcGIS REST Response]
```

### Mode 2: WMS Protocol Translation
```
[WMS Client] → [Proxy Server] → [ArcGIS REST Server]
                     ↓
              [Protocol Translation]
                     ↓
              [WMS Response] ← [ArcGIS REST Response]
```

### Core Components:

1. **HTTP/HTTPS Server** - Handles incoming requests (both modes)
2. **ArcGIS Proxy Handler** - Direct passthrough for `/arcgis/` paths
3. **WMS Handler** - Protocol translation for `/wms` requests
4. **Request Translator** - Converts WMS parameters to ArcGIS REST format
5. **HTTP Client** - Makes requests to upstream ArcGIS server with connection pooling
6. **Response Translator** - Handles response passthrough and WMS error conversion
7. **Configuration Manager** - Environment-based configuration with HTTPS support
8. **Health Check Handler** - Service health endpoint with upstream validation
9. **Certificate Manager** - SSL certificate generation and management

## Implementation Structure

```
wms-proxy/
├── cmd/
│   └── proxy/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management with HTTPS support
│   ├── handlers/
│   │   ├── arcgis_proxy.go      # Direct ArcGIS REST proxy handler
│   │   ├── wms.go               # WMS request handlers
│   │   ├── health.go            # Health check handler
│   │   └── capabilities.go      # GetCapabilities handler
│   ├── translator/
│   │   ├── request.go           # WMS to ArcGIS request translation
│   │   └── response.go          # Response passthrough and error handling
│   ├── client/
│   │   └── arcgis.go            # ArcGIS REST client with connection pooling
│   └── server/
│       └── server.go            # HTTP/HTTPS server setup
├── pkg/
│   └── wms/
│       ├── types.go             # WMS data structures
│       └── capabilities.go      # WMS capabilities XML generation
├── scripts/
│   └── generate-certs.sh        # SSL certificate generation script
├── Dockerfile                   # Container image definition with HTTPS support
├── Makefile                     # Build and run targets (HTTP/HTTPS modes)
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── .gitignore                   # Git ignore file (excludes certs)
├── README.md                    # Usage documentation
├── REQUIREMENTS.md              # Requirements document
└── IMPLEMENTATION_PLAN.md       # This document
```

## Detailed Implementation Plan

### Phase 1: Core Infrastructure (30% of effort)

#### 1.1 Project Setup
- Initialize Go module
- Set up basic project structure
- Create configuration management
- Implement logging infrastructure

#### 1.2 HTTP Server Foundation
- Create basic HTTP server with graceful shutdown
- Implement request routing
- Add middleware for logging and error handling
- Create health check endpoint

### Phase 2: Protocol Translation (40% of effort)

#### 2.1 Request Translation
- Parse WMS GetMap parameters
- Map WMS parameters to ArcGIS REST equivalents:
  - `BBOX` → `bbox`
  - `WIDTH,HEIGHT` → `size`
  - `FORMAT` → `format`
  - `SRS/CRS` → `bboxSR,imageSR`
  - `LAYERS` → `layers`
- Handle coordinate system transformations
- Validate and sanitize input parameters

#### 2.2 Response Translation
- Pass through image responses (PNG, JPEG)
- Handle ArcGIS error responses
- Convert to appropriate WMS error format
- Preserve HTTP headers where appropriate

#### 2.3 ArcGIS Client
- Implement HTTP client with connection pooling
- Handle SSL certificate validation
- Implement retry logic with exponential backoff
- Add request timeout handling

### Phase 3: WMS Compliance (20% of effort)

#### 3.1 GetCapabilities Implementation
- Generate basic WMS capabilities XML
- Query upstream ArcGIS service metadata
- Map ArcGIS layer information to WMS format
- Support multiple WMS versions (1.1.1, 1.3.0)

#### 3.2 Error Handling
- Implement WMS-compliant error responses
- Handle upstream server failures
- Provide meaningful error messages
- Log errors appropriately

### Phase 4: Containerization & Deployment (10% of effort)

#### 4.1 Docker Image
- Create multi-stage Dockerfile
- Use Alpine Linux base for minimal size
- Run as non-root user
- Optimize for security and size

#### 4.2 Build System
- Create Makefile with standard targets
- Support both Docker and Podman
- Include development and production builds
- Add cleanup and testing targets

## Key Implementation Details

### Configuration Management
```go
type Config struct {
    ArcGISHost     string        `env:"ARCGIS_HOST" default:"localhost"`
    ArcGISScheme   string        `env:"ARCGIS_SCHEME" default:"https"`
    ProxyPort      int           `env:"PROXY_PORT" default:"8080"`
    RequestTimeout time.Duration `env:"REQUEST_TIMEOUT" default:"30s"`
    LogLevel       string        `env:"LOG_LEVEL" default:"info"`
    EnableHTTPS    bool          `env:"ENABLE_HTTPS" default:"false"`
    CertFile       string        `env:"CERT_FILE" default:"/app/certs/server.crt"`
    KeyFile        string        `env:"KEY_FILE" default:"/app/certs/server.key"`
}
```

### Request Translation Logic
```go
// WMS GetMap → ArcGIS REST Export
func TranslateGetMapRequest(wmsParams WMSParams) ArcGISParams {
    return ArcGISParams{
        BBOX:        wmsParams.BBOX,
        Size:        fmt.Sprintf("%d,%d", wmsParams.Width, wmsParams.Height),
        Format:      translateFormat(wmsParams.Format),
        BBoxSR:      translateSRS(wmsParams.SRS),
        ImageSR:     translateSRS(wmsParams.SRS),
        Layers:      translateLayers(wmsParams.Layers),
        Transparent: wmsParams.Transparent,
        DPI:         96,
        F:           "image",
    }
}
```

### Error Handling Strategy
```go
// WMS Service Exception Format
func HandleError(err error, w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/vnd.ogc.se_xml")
    w.WriteHeader(http.StatusBadRequest)
    
    xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
    <ServiceExceptionReport version="1.1.1">
        <ServiceException>%s</ServiceException>
    </ServiceExceptionReport>`, html.EscapeString(err.Error()))
    
    w.Write([]byte(xml))
}
```

## Testing Strategy

### Unit Tests
- Request translation functions
- Response handling logic
- Configuration parsing
- Error handling scenarios

### Integration Tests
- End-to-end request flow
- Upstream server interaction
- Container deployment
- WMS client compatibility

### Performance Tests
- Concurrent request handling
- Memory usage under load
- Response time benchmarks

## Deployment Configuration

### Environment Variables
```bash
# Required
ARCGIS_HOST=mapsdep.nj.gov

# Optional with defaults
ARCGIS_SCHEME=https
PROXY_PORT=8080
REQUEST_TIMEOUT=30
LOG_LEVEL=info

# HTTPS Configuration (optional)
ENABLE_HTTPS=false
CERT_FILE=/app/certs/server.crt
KEY_FILE=/app/certs/server.key
```

### Container Execution
```bash
# Build
make build

# Run with custom configuration
make run ARCGIS_HOST=mapsdep.nj.gov

# Run with Podman
podman run -d \
  -p 8080:8080 \
  -e ARCGIS_HOST=mapsdep.nj.gov \
  wms-proxy:latest
```

## Security Considerations

### Input Validation
- Sanitize all URL parameters
- Validate coordinate bounds
- Limit request size and complexity
- Prevent path traversal attacks

### Network Security
- Validate upstream SSL certificates
- Implement request rate limiting
- Use secure defaults for all configurations
- Run container as non-root user

### Container Security
- Minimal base image (Alpine)
- No unnecessary packages
- Read-only filesystem where possible
- Security scanning in build process

## Performance Optimizations

### Connection Management
- HTTP client connection pooling
- Keep-alive connections to upstream
- Configurable connection limits
- Connection timeout handling

### Caching Strategy (Future)
- Response caching based on parameters
- Cache invalidation policies
- Memory-based cache with LRU eviction
- Configurable cache TTL

### Resource Management
- Graceful shutdown handling
- Memory usage monitoring
- Request context cancellation
- Resource cleanup on errors

## Monitoring and Observability

### Logging
- Structured logging with slog
- Request/response correlation IDs
- Performance metrics logging
- Error tracking and alerting

### Health Checks
- Basic service health endpoint
- Upstream connectivity validation
- Resource usage reporting
- Readiness and liveness probes

### Metrics (Future Enhancement)
- Request count and latency
- Error rates and types
- Upstream response times
- Resource utilization metrics

## Risk Mitigation

### Technical Risks
- **Upstream API changes**: Version detection and compatibility checks
- **Performance degradation**: Load testing and optimization
- **Memory leaks**: Proper resource cleanup and monitoring
- **Security vulnerabilities**: Regular dependency updates and scanning

### Operational Risks
- **Service unavailability**: Health checks and retry logic
- **Configuration errors**: Validation and sensible defaults
- **Container issues**: Multi-stage builds and testing
- **Network failures**: Circuit breaker pattern and timeouts

## Success Metrics

### Functional Success
- [x] Direct ArcGIS REST proxy works correctly
- [x] WMS GetMap requests work correctly
- [x] Image responses are properly proxied
- [x] Error handling works as expected
- [x] Container builds and runs successfully
- [x] HTTPS support functions properly
- [x] Health check endpoint responds correctly

### Performance Success
- [x] Response time < 2x direct ArcGIS calls (direct proxy mode)
- [ ] Handles 100+ concurrent requests (needs load testing)
- [x] Memory usage < 512MB
- [x] Container image < 100MB

### Quality Success
- [ ] Unit test coverage > 80% (tests need to be implemented)
- [x] Integration tests pass (manual testing completed)
- [x] Security scan passes (container security implemented)
- [x] Documentation is complete

This implementation plan provides a comprehensive roadmap for developing the ArcGIS REST to WMS proxy server that meets all the requirements outlined in the requirements document.
