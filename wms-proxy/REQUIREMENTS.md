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

#### 1.4 Coordinate Transformation System
- **REQ-014**: The proxy MUST automatically transform coordinates between different spatial reference systems
- **REQ-015**: The proxy MUST support EPSG:3857 (Web Mercator), EPSG:3424 (NAD83 New Jersey), and EPSG:4326 (WGS84) coordinate systems
- **REQ-016**: The proxy MUST dynamically detect backend ArcGIS server coordinate system requirements
- **REQ-017**: The proxy MUST cache backend spatial reference metadata for performance optimization
- **REQ-018**: The proxy MUST handle ESRI WKID to EPSG code mappings (e.g., EPSG:102711 → EPSG:3424)
- **REQ-019**: The proxy MUST only transform coordinates when source ≠ target coordinate system
- **REQ-020**: The proxy MUST provide graceful fallback behavior when coordinate transformation fails

#### 1.5 HTTPS Support
- **REQ-021**: The proxy MUST support HTTPS connections with SSL/TLS certificates
- **REQ-022**: The proxy MUST provide certificate generation for development/testing
- **REQ-023**: The proxy MUST validate SSL certificates for upstream connections
- **REQ-024**: The proxy MUST support both HTTP and HTTPS modes simultaneously

### 2. Configuration Requirements

#### 2.1 Runtime Configuration
- **REQ-025**: The proxy MUST accept target host configuration via environment variable
- **REQ-026**: The proxy MUST accept listening port configuration via environment variable
- **REQ-027**: The proxy MUST support HTTPS upstream servers
- **REQ-028**: The proxy MUST allow configuration of request timeout values
- **REQ-029**: The proxy MUST accept SSL certificate file paths via environment variables
- **REQ-030**: The proxy MUST accept log level configuration via environment variable

#### 2.2 Service Discovery
- **REQ-031**: The proxy MUST support multiple upstream ArcGIS servers (future enhancement)
- **REQ-032**: The proxy MUST validate upstream server connectivity on startup
- **REQ-033**: The proxy MUST automatically query backend service metadata for spatial reference detection
- **REQ-034**: The proxy MUST cache backend spatial reference information with configurable TTL

### 3. WMS Compliance Requirements

#### 3.1 Standard Operations
- **REQ-035**: The proxy MUST support WMS GetMap operations
- **REQ-036**: The proxy MUST support WMS GetCapabilities operations (basic implementation)
- **REQ-037**: The proxy MUST return appropriate WMS-compliant HTTP headers
- **REQ-038**: The proxy MUST handle WMS version parameter (1.1.1, 1.3.0)

#### 3.2 Parameter Mapping
- **REQ-039**: The proxy MUST map WMS BBOX to ArcGIS REST bbox parameter with automatic coordinate transformation
- **REQ-040**: The proxy MUST map WMS WIDTH/HEIGHT to ArcGIS REST size parameter
- **REQ-041**: The proxy MUST map WMS FORMAT to ArcGIS REST format parameter
- **REQ-042**: The proxy MUST map WMS SRS/CRS to ArcGIS REST spatial reference parameters with dynamic backend detection

### 4. Performance Requirements

#### 4.1 Throughput
- **REQ-043**: The proxy MUST handle at least 100 concurrent requests
- **REQ-044**: The proxy MUST maintain response times within 2x of direct ArcGIS REST calls
- **REQ-045**: The proxy MUST implement connection pooling for upstream requests

#### 4.2 Coordinate Transformation Performance
- **REQ-046**: The proxy MUST perform coordinate transformations in under 10 microseconds per bounding box
- **REQ-047**: The proxy MUST cache backend spatial reference metadata with 15-minute TTL to minimize queries
- **REQ-048**: The proxy MUST reduce backend metadata queries by at least 95% through intelligent caching
- **REQ-049**: The proxy MUST add less than 1ms overhead to request processing for coordinate transformations

#### 4.3 Resource Usage
- **REQ-050**: The proxy MUST run within 512MB memory limit in container
- **REQ-051**: The proxy MUST support graceful shutdown on SIGTERM

### 5. Containerization Requirements

#### 5.1 Docker/Podman Support
- **REQ-052**: The proxy MUST run in a container image based on Alpine Linux or similar minimal base
- **REQ-053**: The container MUST expose the service port via EXPOSE directive
- **REQ-054**: The container MUST run as non-root user for security
- **REQ-055**: The container image MUST be under 100MB when compressed

#### 5.2 Build System
- **REQ-056**: A Makefile MUST provide 'build' target for creating container image
- **REQ-057**: A Makefile MUST provide 'run' target for running the container
- **REQ-058**: A Makefile MUST provide 'clean' target for cleanup
- **REQ-059**: The build system MUST work on Fedora Linux 41+

### 6. Security Requirements

#### 6.1 Network Security
- **REQ-060**: The proxy MUST validate upstream SSL certificates
- **REQ-061**: The proxy MUST sanitize input parameters to prevent injection attacks
- **REQ-062**: The proxy MUST implement request rate limiting (basic)
- **REQ-063**: The proxy MUST validate coordinate bounds to prevent malicious coordinate transformation requests

#### 6.2 Container Security
- **REQ-064**: The container MUST run with minimal privileges
- **REQ-065**: The container MUST not include unnecessary packages or tools

### 7. Monitoring and Logging Requirements

#### 7.1 Logging
- **REQ-066**: The proxy MUST log all incoming requests with timestamp
- **REQ-067**: The proxy MUST log upstream request/response status
- **REQ-068**: The proxy MUST log errors with appropriate detail level
- **REQ-069**: Logs MUST be written to stdout for container compatibility
- **REQ-070**: The proxy MUST log coordinate transformation operations with source/target CRS and performance metrics
- **REQ-071**: The proxy MUST log backend spatial reference detection events and cache operations

#### 7.2 Health Checks
- **REQ-072**: The proxy MUST provide a health check endpoint (/health)
- **REQ-073**: The health check MUST verify upstream server connectivity

## Non-Functional Requirements

### 8. Reliability
- **REQ-074**: The proxy MUST handle upstream server failures gracefully
- **REQ-075**: The proxy MUST retry failed upstream requests (configurable attempts)
- **REQ-076**: The proxy MUST continue operating if upstream server is temporarily unavailable
- **REQ-077**: The proxy MUST provide graceful fallback when coordinate transformation fails
- **REQ-078**: The proxy MUST handle malformed backend spatial reference metadata gracefully

### 9. Maintainability
- **REQ-079**: The code MUST include comprehensive error handling
- **REQ-080**: The code MUST be documented with inline comments
- **REQ-081**: The project MUST include a README with usage instructions
- **REQ-082**: The code MUST include comprehensive unit tests for coordinate transformation functions
- **REQ-083**: The project MUST include integration tests for end-to-end coordinate transformation workflows

### 10. Compatibility
- **REQ-084**: The proxy MUST work with standard WMS clients (QGIS, OpenLayers, etc.)
- **REQ-085**: The proxy MUST support ArcGIS REST API versions 10.x and newer
- **REQ-086**: The container MUST run on both Docker and Podman
- **REQ-087**: The proxy MUST work with any ArcGIS backend service regardless of its coordinate system
- **REQ-088**: The proxy MUST handle both modern EPSG codes and legacy ESRI WKID codes

## Example Usage Scenarios

### Scenario 1: Direct ArcGIS REST Proxy (Primary Mode)
```
Client Request (Direct ArcGIS REST):
GET /arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=EPSG:3857&imageSR=EPSG:3857&size=256,256&f=image&layers=show:17

Proxy Action:
Direct passthrough to upstream ArcGIS server with same URL path and parameters
```

### Scenario 2: WMS Protocol Translation with Coordinate Transformation
```
Client Request (WMS with Web Mercator coordinates):
GET /wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&LAYERS=17&STYLES=&FORMAT=image/png&BGCOLOR=0xFFFFFF&TRANSPARENT=TRUE&SRS=EPSG:3857&BBOX=-8238310.24,4969803.4,-8238016.75,4970096.9&WIDTH=256&HEIGHT=256

Proxy Actions:
1. Detects backend expects EPSG:3424 (via service metadata query)
2. Transforms coordinates: EPSG:3857 → EPSG:3424
3. Forwards transformed request to backend

Proxy Translation (ArcGIS REST with transformed coordinates):
GET /arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=629066.028855,684288.229396,629792.766201,685020.398937&bboxSR=EPSG:3424&imageSR=EPSG:3424&size=256,256&f=image&layers=show:17
```

### Scenario 2a: Automatic Coordinate System Detection
```
Client Request (Direct ArcGIS REST with WGS84 coordinates):
GET /arcgis/rest/services/Features/Environmental_admin/MapServer/export?bbox=-74.006000,40.710974,-74.003364,40.712972&bboxSR=EPSG:4326&imageSR=EPSG:4326&size=256,256&f=image&layers=show:17

Proxy Actions:
1. Queries backend service metadata: GET /arcgis/rest/services/Features/Environmental_admin/MapServer?f=json
2. Detects backend spatial reference: {"spatialReference": {"wkid": 102711, "latestWkid": 3424}}
3. Maps EPSG:102711 → EPSG:3424 and transforms coordinates: EPSG:4326 → EPSG:3424
4. Caches backend SR for 15 minutes

Final Backend Request:
GET /arcgis/rest/services/Features/Environmental_admin/MapServer/export?bbox=629066.039508,684288.262640,629792.648846,685020.247912&bboxSR=EPSG:3424&imageSR=EPSG:4326&size=256,256&f=image&layers=show:17
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

### Scenario 5: Universal Backend Compatibility Testing
```
Test Script Usage (covers entire New Jersey state):
./scripts/run-proxy-test.sh --map deer --sr 3857  # Web Mercator coordinates
./scripts/run-proxy-test.sh --map deer --sr 3424  # New Jersey State Plane coordinates  
./scripts/run-proxy-test.sh --map deer --sr 4326  # WGS84 Geographic coordinates

Each test uses the same geographic area (entire NJ state) but in different coordinate systems:
- EPSG:3857: -8411257.76,4711437.70,-8225840.11,5065205.28
- EPSG:3424: 190699.22,36416.14,658483.94,919995.68
- EPSG:4326: -75.559614,38.928519,-73.893979,41.357423

Proxy automatically transforms all coordinate systems to backend's expected format.
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
3. **Coordinate transformation works automatically** between EPSG:3857, EPSG:3424, and EPSG:4326
4. **Dynamic backend detection** automatically determines backend coordinate system requirements
5. **Universal backend compatibility** works with any ArcGIS service regardless of coordinate system
6. **Performance requirements met** with <10μs coordinate transformations and 95% cache hit rate
7. The container builds and runs successfully on Fedora Linux 41+
8. The Makefile targets work as specified
9. Performance meets the stated requirements
10. All functional requirements are implemented and tested
11. **Comprehensive test coverage** with 18+ unit tests achieving 100% pass rate
12. **End-to-end testing** validates coordinate transformation across entire New Jersey state

## Future Enhancements (Out of Scope)

- Support for WMS GetFeatureInfo operations
- ~~Caching layer for improved performance~~ ✅ **IMPLEMENTED** - Backend SR metadata caching with 15-minute TTL
- Support for multiple simultaneous upstream servers
- Advanced authentication/authorization
- Metrics and monitoring integration
- Support for vector tile formats
- Additional coordinate system support beyond EPSG:3857, EPSG:3424, EPSG:4326
- Advanced coordinate transformation algorithms for higher precision
- Distributed caching (Redis, database) for multi-instance deployments
