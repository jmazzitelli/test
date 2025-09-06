# ArcGIS REST to WMS Proxy Server - Requirements Document

## Overview

This document outlines the requirements for developing a proxy server that acts as an intermediary between clients and ArcGIS REST API servers. The proxy supports two operational modes:

1. **Direct ArcGIS REST Proxy**: Transparent proxy that forwards ArcGIS REST requests directly
2. **WMS Protocol Translation**: Converts WMS requests to ArcGIS REST format for WMS clients

This dual-mode approach provides maximum flexibility for different client types.

## Functional Requirements

### 1. Core Proxy Functionality

#### 1.1 Request Forwarding
- **REQ-001**: The proxy MUST accept HTTP and HTTPS requests on configurable ports (default: 8080 HTTP, 8443 HTTPS)
- **REQ-002**: The proxy MUST forward requests to a configurable ArcGIS REST server host
- **REQ-003**: The proxy MUST preserve all query parameters from the original request when forwarding
- **REQ-004**: The proxy MUST handle URL path mapping between WMS and ArcGIS REST endpoints
- **REQ-005**: The proxy MUST support direct ArcGIS REST API passthrough for `/arcgis/` paths

#### 1.2 Response Translation
- **REQ-006**: The proxy MUST convert ArcGIS REST API responses to WMS-compliant format (WMS mode only)
- **REQ-007**: The proxy MUST handle image format responses (PNG, JPEG, etc.)
- **REQ-008**: The proxy MUST preserve image data integrity during translation
- **REQ-009**: The proxy MUST handle error responses and translate them to WMS error format

#### 1.3 Protocol Translation
- **REQ-010**: The proxy MUST support ArcGIS REST MapServer export operations
- **REQ-011**: The proxy MUST translate coordinate reference systems (CRS) between formats (WMS mode)
- **REQ-012**: The proxy MUST handle bounding box (BBOX) parameter translation (WMS mode)
- **REQ-013**: The proxy MUST support common image formats (PNG32, PNG, JPEG)

#### 1.4 HTTPS Support
- **REQ-014**: The proxy MUST support HTTPS connections with SSL/TLS certificates
- **REQ-015**: The proxy MUST provide certificate generation for development/testing
- **REQ-016**: The proxy MUST validate SSL certificates for upstream connections
- **REQ-017**: The proxy MUST support both HTTP and HTTPS modes simultaneously

### 2. Configuration Requirements

#### 2.1 Runtime Configuration
- **REQ-018**: The proxy MUST accept target host configuration via environment variable
- **REQ-019**: The proxy MUST accept listening port configuration via environment variable
- **REQ-020**: The proxy MUST support HTTPS upstream servers
- **REQ-021**: The proxy MUST allow configuration of request timeout values
- **REQ-022**: The proxy MUST accept SSL certificate file paths via environment variables
- **REQ-023**: The proxy MUST accept log level configuration via environment variable

#### 2.2 Service Discovery
- **REQ-024**: The proxy MUST support multiple upstream ArcGIS servers (future enhancement)
- **REQ-025**: The proxy MUST validate upstream server connectivity on startup

### 3. WMS Compliance Requirements

#### 3.1 Standard Operations
- **REQ-026**: The proxy MUST support WMS GetMap operations
- **REQ-027**: The proxy MUST support WMS GetCapabilities operations (basic implementation)
- **REQ-028**: The proxy MUST return appropriate WMS-compliant HTTP headers
- **REQ-029**: The proxy MUST handle WMS version parameter (1.1.1, 1.3.0)

#### 3.2 Parameter Mapping
- **REQ-030**: The proxy MUST map WMS BBOX to ArcGIS REST bbox parameter
- **REQ-031**: The proxy MUST map WMS WIDTH/HEIGHT to ArcGIS REST size parameter
- **REQ-032**: The proxy MUST map WMS FORMAT to ArcGIS REST format parameter
- **REQ-033**: The proxy MUST map WMS SRS/CRS to ArcGIS REST spatial reference parameters

### 4. Performance Requirements

#### 4.1 Throughput
- **REQ-027**: The proxy MUST handle at least 100 concurrent requests
- **REQ-028**: The proxy MUST maintain response times within 2x of direct ArcGIS REST calls
- **REQ-029**: The proxy MUST implement connection pooling for upstream requests

#### 4.2 Resource Usage
- **REQ-030**: The proxy MUST run within 512MB memory limit in container
- **REQ-031**: The proxy MUST support graceful shutdown on SIGTERM

### 5. Containerization Requirements

#### 5.1 Docker/Podman Support
- **REQ-032**: The proxy MUST run in a container image based on Alpine Linux or similar minimal base
- **REQ-033**: The container MUST expose the service port via EXPOSE directive
- **REQ-034**: The container MUST run as non-root user for security
- **REQ-035**: The container image MUST be under 100MB when compressed

#### 5.2 Build System
- **REQ-036**: A Makefile MUST provide 'build' target for creating container image
- **REQ-037**: A Makefile MUST provide 'run' target for running the container
- **REQ-038**: A Makefile MUST provide 'clean' target for cleanup
- **REQ-039**: The build system MUST work on Fedora Linux 41+

### 6. Security Requirements

#### 6.1 Network Security
- **REQ-040**: The proxy MUST validate upstream SSL certificates
- **REQ-041**: The proxy MUST sanitize input parameters to prevent injection attacks
- **REQ-042**: The proxy MUST implement request rate limiting (basic)

#### 6.2 Container Security
- **REQ-043**: The container MUST run with minimal privileges
- **REQ-044**: The container MUST not include unnecessary packages or tools

### 7. Monitoring and Logging Requirements

#### 7.1 Logging
- **REQ-045**: The proxy MUST log all incoming requests with timestamp
- **REQ-046**: The proxy MUST log upstream request/response status
- **REQ-047**: The proxy MUST log errors with appropriate detail level
- **REQ-048**: Logs MUST be written to stdout for container compatibility

#### 7.2 Health Checks
- **REQ-049**: The proxy MUST provide a health check endpoint (/health)
- **REQ-050**: The health check MUST verify upstream server connectivity

## Non-Functional Requirements

### 8. Reliability
- **REQ-051**: The proxy MUST handle upstream server failures gracefully
- **REQ-052**: The proxy MUST retry failed upstream requests (configurable attempts)
- **REQ-053**: The proxy MUST continue operating if upstream server is temporarily unavailable

### 9. Maintainability
- **REQ-054**: The code MUST include comprehensive error handling
- **REQ-055**: The code MUST be documented with inline comments
- **REQ-056**: The project MUST include a README with usage instructions

### 10. Compatibility
- **REQ-057**: The proxy MUST work with standard WMS clients (QGIS, OpenLayers, etc.)
- **REQ-058**: The proxy MUST support ArcGIS REST API versions 10.x and newer
- **REQ-059**: The container MUST run on both Docker and Podman

## Example Usage Scenarios

### Scenario 1: Direct ArcGIS REST Proxy (Primary Mode)
```
Client Request (Direct ArcGIS REST):
GET /arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=EPSG:3857&imageSR=EPSG:3857&size=256,256&f=image&layers=show:17

Proxy Action:
Direct passthrough to upstream ArcGIS server with same URL path and parameters
```

### Scenario 2: WMS Protocol Translation
```
Client Request (WMS):
GET /wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&LAYERS=17&STYLES=&FORMAT=image/png&BGCOLOR=0xFFFFFF&TRANSPARENT=TRUE&SRS=EPSG:3857&BBOX=-8238310.24,4969803.4,-8238016.75,4970096.9&WIDTH=256&HEIGHT=256

Proxy Translation (ArcGIS REST):
GET /arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=EPSG:3857&imageSR=EPSG:3857&size=256,256&f=image&layers=show:17
```

### Scenario 3: WMS GetCapabilities Request
```
Client Request:
GET /wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetCapabilities

Proxy Response:
Returns basic WMS capabilities XML describing available layers from upstream ArcGIS service
```

### Scenario 4: HTTPS Support
```
Client Request (HTTPS):
GET https://localhost:8443/arcgis/rest/services/Features/Environmental_admin/MapServer/export?...

Proxy Action:
Accepts HTTPS connection, forwards to upstream server via HTTPS, returns response over HTTPS
```

## Configuration Examples

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
ENABLE_HTTPS=true
CERT_FILE=/app/certs/server.crt
KEY_FILE=/app/certs/server.key
```

### Container Run Examples
```bash
# HTTP Mode
podman run -d \
  -p 8080:8080 \
  -e ARCGIS_HOST=mapsdep.nj.gov \
  -e ARCGIS_SCHEME=https \
  wms-proxy:latest

# HTTPS Mode
podman run -d \
  -p 8080:8080 -p 8443:8443 \
  -v ./certs:/app/certs:Z \
  -e ARCGIS_HOST=mapsdep.nj.gov \
  -e ARCGIS_SCHEME=https \
  -e ENABLE_HTTPS=true \
  -e PROXY_PORT=8443 \
  wms-proxy:latest
```

## Success Criteria

The implementation will be considered successful when:

1. A WMS client can successfully retrieve map tiles through the proxy
2. The proxy correctly translates between WMS and ArcGIS REST protocols
3. The container builds and runs successfully on Fedora Linux 41+
4. The Makefile targets work as specified
5. Performance meets the stated requirements
6. All functional requirements are implemented and tested

## Future Enhancements (Out of Scope)

- Support for WMS GetFeatureInfo operations
- Caching layer for improved performance
- Support for multiple simultaneous upstream servers
- Advanced authentication/authorization
- Metrics and monitoring integration
- Support for vector tile formats
