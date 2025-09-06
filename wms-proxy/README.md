# WMS Proxy Server

A dual-mode proxy server that provides both direct ArcGIS REST API passthrough and WMS (Web Map Service) protocol translation. This allows both existing ArcGIS clients and WMS clients to access ArcGIS REST services seamlessly.

## Features

- **Dual Mode Operation**: Direct ArcGIS REST proxy + WMS protocol translation
- **Direct ArcGIS REST Proxy**: Transparent passthrough for existing ArcGIS clients
- **WMS Protocol Translation**: Converts WMS GetMap requests to ArcGIS REST export requests
- **HTTPS Support**: Full SSL/TLS support with certificate generation
- **Image Passthrough**: Efficiently proxies image responses (PNG, JPEG, GIF)
- **WMS Compliance**: Supports basic WMS operations (GetMap, GetCapabilities)
- **Containerized**: Runs in Docker/Podman containers with multi-arch support
- **Health Monitoring**: Built-in health check endpoint with upstream validation
- **Structured Logging**: JSON-based logging with configurable levels
- **Connection Pooling**: Efficient upstream connection management
- **Configurable**: Environment variable based configuration
- **Security**: Runs as non-root user, input validation, SSL certificate validation

## Quick Start

### Using Podman/Docker

1. **Build and run with default settings:**
   ```bash
   make run ARCGIS_HOST=mapsdep.nj.gov
   ```

2. **Test the proxy:**
   ```bash
   # Health check
   curl http://localhost:8080/health
   
   # GetCapabilities
   curl "http://localhost:8080/wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetCapabilities"
   
   # GetMap example
   curl "http://localhost:8080/wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&LAYERS=17&FORMAT=image/png&BBOX=-8238310.24,4969803.4,-8238016.75,4970096.9&WIDTH=256&HEIGHT=256&SRS=EPSG:3857" -o test.png
   ```

### Using Local Development

1. **Build and run locally (HTTP):**
   ```bash
   make dev-run ARCGIS_HOST=mapsdep.nj.gov
   ```

2. **Build and run locally (HTTPS):**
   ```bash
   make dev-run-https ARCGIS_HOST=mapsdep.nj.gov
   ```

### Using HTTPS

1. **Run with HTTPS support:**
   ```bash
   make run-https ARCGIS_HOST=mapsdep.nj.gov
   ```

2. **Test HTTPS proxy:**
   ```bash
   # Health check (accept self-signed certificate)
   curl -k https://localhost:8443/health
   
   # GetMap example (HTTPS)
   curl -k "https://localhost:8443/arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=EPSG:3857&imageSR=EPSG:3857&size=256,256&f=image&layers=show:17" -o test.png
   ```

**Note:** The `-k` flag is needed with curl to accept self-signed certificates. For production, use certificates from a trusted CA.

## Configuration

The proxy is configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `ARCGIS_HOST` | Target ArcGIS server hostname | `localhost` |
| `ARCGIS_SCHEME` | Protocol for ArcGIS server (http/https) | `https` |
| `ARCGIS_SERVICE` | ArcGIS service path for WMS translation | `/arcgis/rest/services/Features/Environmental_admin/MapServer/export` |
| `PROXY_PORT` | Port for proxy to listen on | `8080` |
| `REQUEST_TIMEOUT` | Timeout for upstream requests (seconds) | `30` |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | `info` |
| `ENABLE_HTTPS` | Enable HTTPS support (true/false) | `false` |
| `CERT_FILE` | Path to SSL certificate file | `/app/certs/server.crt` |
| `KEY_FILE` | Path to SSL private key file | `/app/certs/server.key` |

## Makefile Targets

### Container Operations
- `make build` - Build container image
- `make run` - Build and run container (HTTP)
- `make run-https` - Build and run container with HTTPS
- `make stop` - Stop and remove container
- `make clean` - Remove container and image
- `make logs` - Show container logs
- `make health` - Check container health
- `make gen-certs` - Generate self-signed SSL certificates

### Development
- `make dev-build` - Build Go binary locally
- `make dev-run` - Run Go binary locally (HTTP)
- `make dev-run-https` - Run Go binary locally with HTTPS
- `make test` - Run tests

### Docker Support
- `make docker-build` - Build with Docker instead of Podman
- `make docker-run` - Run with Docker
- `make docker-stop` - Stop Docker container
- `make docker-clean` - Clean Docker resources

### Examples
- `make example` - Show example request URLs

## Usage Examples

The proxy supports **two modes** of operation:

### Mode 1: Direct ArcGIS REST Proxy (Recommended)

Use the **exact same ArcGIS REST URLs** as your original request, just change the host to localhost:

**Original ArcGIS URL:**
```
https://mapsdep.nj.gov/arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=EPSG:3857&imageSR=EPSG:3857&size=256,256&f=image&layers=show:17
```

**Proxy URL (just change host):**
```bash
curl "http://localhost:8080/arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=EPSG:3857&imageSR=EPSG:3857&size=256,256&f=image&layers=show:17" -o map.png
```

**This is what you want!** Your client can use the exact same URL format, just pointing to `localhost:8080` instead of `mapsdep.nj.gov`.

### Mode 2: WMS Client Configuration (For WMS Clients)

Configure your WMS client to use the proxy:

**Base URL:** `http://localhost:8080/wms`

**Example Layer Configuration:**
```
Service: WMS
Version: 1.1.1
URL: http://localhost:8080/wms
Layers: 17
Format: image/png
SRS: EPSG:3857
```

**Manual WMS Requests:**

**GetCapabilities:**
```bash
curl "http://localhost:8080/wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetCapabilities"
```

**GetMap:**
```bash
curl "http://localhost:8080/wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&LAYERS=17&STYLES=&FORMAT=image/png&BGCOLOR=0xFFFFFF&TRANSPARENT=TRUE&SRS=EPSG:3857&BBOX=-8238310.24,4969803.4,-8238016.75,4970096.9&WIDTH=256&HEIGHT=256" -o map.png
```

### QGIS Integration

1. Add a new WMS layer in QGIS
2. Use URL: `http://localhost:8080/wms`
3. Click "Connect" to load available layers
4. Select and add desired layers

## Architecture

```
[WMS Client] → [Proxy Server] → [ArcGIS REST Server]
                     ↓
              [Protocol Translation]
                     ↓
              [WMS Response] ← [ArcGIS REST Response]
```

### Request Translation

The proxy translates WMS parameters to ArcGIS REST format:

| WMS Parameter | ArcGIS Parameter | Notes |
|---------------|------------------|-------|
| `BBOX` | `bbox` | Direct mapping |
| `WIDTH,HEIGHT` | `size` | Combined as "width,height" |
| `FORMAT` | `format` | Translated (png→png32, etc.) |
| `SRS/CRS` | `bboxSR,imageSR` | Spatial reference system |
| `LAYERS` | `layers` | Converted to "show:layerId" format |
| `TRANSPARENT` | `transparent` | Boolean conversion |

### Response Handling

- **Images**: Passed through directly with appropriate headers
- **Errors**: Converted to WMS-compliant error XML
- **Capabilities**: Generated basic WMS capabilities XML

## Monitoring

### Health Check

The proxy provides a health check endpoint at `/health`:

```bash
curl http://localhost:8080/health
```

Response format:
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "upstream": "ok",
  "message": ""
}
```

### Logging

The proxy uses structured JSON logging. Log levels can be controlled via `LOG_LEVEL` environment variable.

Example log entry:
```json
{
  "time": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "msg": "HTTP request",
  "method": "GET",
  "path": "/wms",
  "status": 200,
  "duration_ms": 150,
  "remote_addr": "192.168.1.100:54321"
}
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Check if container is running: `podman ps`
   - Verify port mapping: `-p 8080:8080`

2. **Upstream Server Error**
   - Check health endpoint: `curl http://localhost:8080/health`
   - Verify ARCGIS_HOST configuration
   - Check network connectivity to ArcGIS server

3. **Invalid Parameters**
   - Ensure required WMS parameters are provided
   - Check parameter format and values
   - Review proxy logs: `make logs`

### Debug Mode

Run with debug logging:
```bash
make dev-run LOG_LEVEL=debug ARCGIS_HOST=mapsdep.nj.gov
```

### Container Debugging

Access container logs:
```bash
make logs
```

Execute commands in container:
```bash
podman exec -it wms-proxy-container /bin/sh
```

## Development

### Project Structure

```
wms-proxy/
├── cmd/proxy/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── translator/      # Protocol translation logic
│   ├── client/          # ArcGIS REST client
│   └── server/          # HTTP server setup
├── pkg/wms/             # WMS data structures
├── Dockerfile           # Container definition
├── Makefile            # Build automation
└── README.md           # This file
```

### Building from Source

1. **Prerequisites:**
   - Go 1.21 or later
   - Podman or Docker

2. **Build:**
   ```bash
   go mod tidy
   go build -o wms-proxy ./cmd/proxy
   ```

3. **Run:**
   ```bash
   ARCGIS_HOST=mapsdep.nj.gov ./wms-proxy
   ```

### Testing

Run tests:
```bash
make test
```

## Security Considerations

- Container runs as non-root user
- Input validation on all parameters
- SSL certificate validation for upstream requests
- No sensitive data in logs
- Minimal container image (Alpine-based)

## Performance

- Connection pooling for upstream requests
- Efficient image streaming (no buffering)
- Configurable timeouts
- Graceful shutdown handling

## License

This project is provided as-is for demonstration purposes.

## Support

For issues and questions:
1. Check the troubleshooting section
2. Review container logs
3. Verify configuration
4. Test health endpoint
