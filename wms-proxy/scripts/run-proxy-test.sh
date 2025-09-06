#!/bin/bash

# Test script for WMS Proxy
# This script starts the proxy, runs a test request, captures logs, and stops the proxy

set -e  # Exit on any error

# Parse command line arguments
MAP_TYPE="waterfowl"  # Default
while [[ $# -gt 0 ]]; do
  case $1 in
    --map)
      MAP_TYPE="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--map waterfowl|deer]"
      exit 1
      ;;
  esac
done

# Validate map type
if [[ "$MAP_TYPE" != "waterfowl" && "$MAP_TYPE" != "deer" ]]; then
  echo "Error: --map must be either 'waterfowl' or 'deer'"
  exit 1
fi

# Change to the parent directory to where this script is
SCRIPT_ROOT="$( cd "$(dirname "$0")" ; pwd -P )"
cd ${SCRIPT_ROOT}/..

echo "=== WMS Proxy Test Script ==="
echo "Map type: $MAP_TYPE"
echo "Starting proxy with HTTPS..."

# Start the proxy using make target
make run-https

echo "Waiting for proxy to start..."
sleep 3

# Test the proxy based on map type
echo "Submitting test request..."
if [[ "$MAP_TYPE" == "waterfowl" ]]; then
  # Waterfowl hunting zones (WMS request)
  curl -k "https://localhost:8443/wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&LAYERS=11&STYLES=&FORMAT=image/png&BGCOLOR=0xFFFFFF&TRANSPARENT=TRUE&SRS=EPSG:3424&BBOX=400000,100000,800000,500000&WIDTH=512&HEIGHT=512" \
       -o test-waterfowl-image.png
  IMAGE_FILE="test-waterfowl-zones.png"
elif [[ "$MAP_TYPE" == "deer" ]]; then
  # Deer management zones (ArcGIS REST request)
  curl -k "https://localhost:8443/arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=100000,0,800000,500000&bboxSR=3424&imageSR=3424&size=512,512&f=image&layers=show:17" \
       -o test-deer-image.png
  IMAGE_FILE="test-deer-zones.png"
fi

echo "Test image saved to: $IMAGE_FILE"

# Capture all proxy logs
echo "Capturing proxy logs..."
podman logs wms-proxy-container > test-proxy-run.log

echo "Logs saved to: test-proxy-run.log"

# Stop the proxy
echo "Stopping proxy..."
make stop

echo "=== Test Complete ==="
echo "Files created:"
echo "  - $IMAGE_FILE (map image)"
echo "  - test-proxy-run.log (proxy logs)"
echo ""
echo "Check the logs to see the REST URL that was generated!"
