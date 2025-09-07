# WMS Proxy Makefile

# Variables
IMAGE_NAME := wms-proxy
IMAGE_TAG := latest
CONTAINER_NAME := wms-proxy-container
PROXY_PORT := 8080
HTTPS_PORT := 8443
ARCGIS_HOST := mapsdep.nj.gov
ARCGIS_SCHEME := https
ENABLE_HTTPS := false

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build       - Build the container image"
	@echo "  run         - Run the container (requires ARCGIS_HOST)"
	@echo "  stop        - Stop the running container"
	@echo "  clean       - Remove container and image"
	@echo "  logs        - Show container logs"
	@echo "  health      - Check container health"
	@echo "  dev-build   - Build Go binary locally"
	@echo "  dev-run     - Run Go binary locally"
	@echo "  test        - Run tests"
	@echo "  gen-certs   - Generate self-signed SSL certificates"
	@echo "  run-https   - Run container with HTTPS enabled"
	@echo ""
	@echo "Environment variables:"
	@echo "  ARCGIS_HOST   - Target ArcGIS server host (default: $(ARCGIS_HOST))"
	@echo "  ARCGIS_SCHEME - Target ArcGIS server scheme (default: $(ARCGIS_SCHEME))"
	@echo "  PROXY_PORT    - HTTP proxy listening port (default: $(PROXY_PORT))"
	@echo "  HTTPS_PORT    - HTTPS proxy listening port (default: $(HTTPS_PORT))"
	@echo "  ENABLE_HTTPS  - Enable HTTPS support (default: $(ENABLE_HTTPS))"

# Build container image
.PHONY: build
build:
	@echo "Building container image..."
	podman build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "Build complete: $(IMAGE_NAME):$(IMAGE_TAG)"

# Run container
.PHONY: run
run: build
	@echo "Starting WMS proxy container..."
	@echo "ArcGIS Host: $(ARCGIS_HOST)"
	@echo "Proxy Port: $(PROXY_PORT)"
	podman run -d \
		--name $(CONTAINER_NAME) \
		-p $(PROXY_PORT):8080 \
		-e ARCGIS_HOST=$(ARCGIS_HOST) \
		-e ARCGIS_SCHEME=$(ARCGIS_SCHEME) \
		-e PROXY_PORT=8080 \
		-e LOG_LEVEL=info \
		$(IMAGE_NAME):$(IMAGE_TAG)
	@echo "Container started. Access proxy at http://localhost:$(PROXY_PORT)"
	@echo "Health check: http://localhost:$(PROXY_PORT)/health"

# Stop container
.PHONY: stop
stop:
	@echo "Stopping container..."
	-podman stop $(CONTAINER_NAME)
	-podman rm $(CONTAINER_NAME)
	@echo "Container stopped and removed"

# Clean up
.PHONY: clean
clean: stop
	@echo "Cleaning up images..."
	-podman rmi $(IMAGE_NAME):$(IMAGE_TAG)
	rm $(CURDIR)/wms-proxy || true
	@echo "Cleaning up test files..."
	@rm $(CURDIR)/test-* || true
	@echo "Cleanup complete"

# Show logs
.PHONY: logs
logs:
	podman logs -f $(CONTAINER_NAME)

# Health check
.PHONY: health
health:
	@echo "Checking container health..."
	@curl -s http://localhost:$(PROXY_PORT)/health | jq . || echo "Health check failed or jq not available"

# Development build (local Go binary)
.PHONY: dev-build
dev-build:
	@echo "Building Go binary locally..."
	go mod tidy
	go build -o wms-proxy ./cmd/proxy
	@echo "Binary built: ./wms-proxy"

# Development run (local Go binary)
.PHONY: dev-run
dev-run: dev-build
	@echo "Running locally with environment variables..."
	ARCGIS_HOST=$(ARCGIS_HOST) \
	ARCGIS_SCHEME=$(ARCGIS_SCHEME) \
	PROXY_PORT=$(PROXY_PORT) \
	LOG_LEVEL=debug \
	./wms-proxy

# Development run with HTTPS (local Go binary)
.PHONY: dev-run-https
dev-run-https: dev-build gen-certs
	@echo "Running locally with HTTPS enabled..."
	ARCGIS_HOST=$(ARCGIS_HOST) \
	ARCGIS_SCHEME=$(ARCGIS_SCHEME) \
	PROXY_PORT=$(HTTPS_PORT) \
	ENABLE_HTTPS=true \
	CERT_FILE=$(PWD)/certs/server.crt \
	KEY_FILE=$(PWD)/certs/server.key \
	LOG_LEVEL=debug \
	./wms-proxy

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Generate SSL certificates
.PHONY: gen-certs
gen-certs:
	@echo "Generating self-signed SSL certificates..."
	./scripts/generate-certs.sh

# Run container with HTTPS enabled
.PHONY: run-https
run-https: build gen-certs
	@echo "Starting WMS proxy container with HTTPS..."
	@echo "ArcGIS Host: $(ARCGIS_HOST)"
	@echo "HTTP Port: $(PROXY_PORT)"
	@echo "HTTPS Port: $(HTTPS_PORT)"
	podman run -d \
		--name $(CONTAINER_NAME) \
		-p $(PROXY_PORT):8080 \
		-p $(HTTPS_PORT):8443 \
		-v $(PWD)/certs:/app/certs:Z \
		-e ARCGIS_HOST=$(ARCGIS_HOST) \
		-e ARCGIS_SCHEME=$(ARCGIS_SCHEME) \
		-e PROXY_PORT=8443 \
		-e ENABLE_HTTPS=true \
		-e CERT_FILE=/app/certs/server.crt \
		-e KEY_FILE=/app/certs/server.key \
		-e LOG_LEVEL=info \
		$(IMAGE_NAME):$(IMAGE_TAG)
	@echo "Container started with HTTPS support"
	@echo "HTTP access: http://localhost:$(PROXY_PORT)"
	@echo "HTTPS access: https://localhost:$(HTTPS_PORT)"
	@echo "Health check: https://localhost:$(HTTPS_PORT)/health"

# Docker variants (for systems that prefer Docker over Podman)
.PHONY: docker-build
docker-build:
	@echo "Building with Docker..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: docker-run
docker-run: docker-build
	@echo "Running with Docker..."
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(PROXY_PORT):8080 \
		-e ARCGIS_HOST=$(ARCGIS_HOST) \
		-e ARCGIS_SCHEME=$(ARCGIS_SCHEME) \
		-e PROXY_PORT=8080 \
		-e LOG_LEVEL=info \
		$(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	-docker stop $(CONTAINER_NAME)
	-docker rm $(CONTAINER_NAME)

.PHONY: docker-clean
docker-clean: docker-stop
	@echo "Cleaning Docker images..."
	-docker rmi $(IMAGE_NAME):$(IMAGE_TAG)

# Example usage targets
.PHONY: example
example:
	@echo "Example request URLs:"
	@echo ""
	@echo "=== Direct ArcGIS REST Proxy (Recommended) ==="
	@echo "Use the same ArcGIS URLs, just change host to localhost:"
	@echo ""
	@echo "Original ArcGIS URL:"
	@echo "  https://$(ARCGIS_HOST)/arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=EPSG:3857&imageSR=EPSG:3857&size=256,256&f=image&layers=show:17"
	@echo ""
	@echo "Proxy URL (HTTP):"
	@echo "  http://localhost:$(PROXY_PORT)/arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=EPSG:3857&imageSR=EPSG:3857&size=256,256&f=image&layers=show:17"
	@echo ""
	@echo "Proxy URL (HTTPS):"
	@echo "  https://localhost:$(HTTPS_PORT)/arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=-8238310.24,4969803.4,-8238016.75,4970096.9&bboxSR=EPSG:3857&imageSR=EPSG:3857&size=256,256&f=image&layers=show:17"
	@echo ""
	@echo "=== WMS Mode (For WMS Clients) ==="
	@echo "GetCapabilities:"
	@echo "  http://localhost:$(PROXY_PORT)/wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetCapabilities"
	@echo ""
	@echo "GetMap (example):"
	@echo "  http://localhost:$(PROXY_PORT)/wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&LAYERS=17&FORMAT=image/png&BBOX=-8238310.24,4969803.4,-8238016.75,4970096.9&WIDTH=256&HEIGHT=256&SRS=EPSG:3857"
	@echo ""
	@echo "=== Health Check ==="
	@echo "  http://localhost:$(PROXY_PORT)/health"
