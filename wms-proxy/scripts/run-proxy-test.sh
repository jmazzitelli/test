#!/bin/bash

# Test script for WMS Proxy
# This script starts the proxy, runs a test request, captures logs, and stops the proxy

set -e  # Exit on any error

# Parse command line arguments
MAP_TYPE="waterfowl"  # Default
SR_TYPE="3424"        # Default coordinate system (New Jersey State Plane)
while [[ $# -gt 0 ]]; do
  case $1 in
    --map)
      MAP_TYPE="$2"
      shift 2
      ;;
    --sr)
      SR_TYPE="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      echo "Usage: $0 [--map waterfowl|deer] [--sr 3857|3424|4326]"
      exit 1
      ;;
  esac
done

# Validate map type
if [[ "$MAP_TYPE" != "waterfowl" && "$MAP_TYPE" != "deer" ]]; then
  echo "Error: --map must be either 'waterfowl' or 'deer'"
  exit 1
fi

# Validate spatial reference system
if [[ "$SR_TYPE" != "3857" && "$SR_TYPE" != "3424" && "$SR_TYPE" != "4326" ]]; then
  echo "Error: --sr must be one of:"
  echo "  3857 - Web Mercator"
  echo "  3424 - New Jersey State Plane"
  echo "  4326 - WGS84 Geographic (lat/lon)"
  exit 1
fi

# Define coordinate sets covering ALL of New Jersey in different coordinate systems
# These coordinates represent the exact same geographic area (entire NJ state) in different projections
# Bounds: West: -75.559614, East: -73.893979, South: 38.928519, North: 41.357423
if [[ "$SR_TYPE" == "3857" ]]; then
  # Web Mercator (EPSG:3857) coordinates - entire New Jersey state
  NJ_BBOX="-8411257.76,4711437.70,-8225840.11,5065205.28"
  WATERFOWL_BBOX="$NJ_BBOX"
  DEER_BBOX="$NJ_BBOX"
elif [[ "$SR_TYPE" == "3424" ]]; then
  # New Jersey State Plane (EPSG:3424) coordinates - entire New Jersey state
  NJ_BBOX="190699.22,36416.14,658483.94,919995.68"
  WATERFOWL_BBOX="$NJ_BBOX"
  DEER_BBOX="$NJ_BBOX"
elif [[ "$SR_TYPE" == "4326" ]]; then
  # WGS84 Geographic (EPSG:4326) coordinates - entire New Jersey state in lat/lon
  NJ_BBOX="-75.559614,38.928519,-73.893979,41.357423"
  WATERFOWL_BBOX="$NJ_BBOX"
  DEER_BBOX="$NJ_BBOX"
fi

# Change to the parent directory to where this script is
SCRIPT_ROOT="$( cd "$(dirname "$0")" ; pwd -P )"
cd ${SCRIPT_ROOT}/..

echo "=== WMS Proxy Test Script ==="
echo "Map type: $MAP_TYPE"
echo "Spatial Reference System: EPSG:$SR_TYPE"
echo "Coverage area: Entire New Jersey state"
echo "Starting proxy with HTTPS..."

# Start the proxy using make target
make run-https

echo "Waiting for proxy to start..."
sleep 3

# Test the proxy based on map type
echo "Submitting test request covering entire New Jersey with coordinates in EPSG:$SR_TYPE format..."
if [[ "$MAP_TYPE" == "waterfowl" ]]; then
  # Waterfowl hunting zones (WMS request) - entire New Jersey coverage
  # The proxy will transform coordinates from SR_TYPE to backend's expected coordinate system
  IMAGE_FILE="test-waterfowl-sr$SR_TYPE.png"
  curl -k "https://localhost:8443/wms?SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&LAYERS=11&STYLES=&FORMAT=image/png&BGCOLOR=0xFFFFFF&TRANSPARENT=TRUE&SRS=EPSG:$SR_TYPE&BBOX=$WATERFOWL_BBOX&WIDTH=512&HEIGHT=512" \
       -o "$IMAGE_FILE"
elif [[ "$MAP_TYPE" == "deer" ]]; then
  # Deer management zones (ArcGIS REST request) - entire New Jersey coverage
  # The proxy will transform coordinates from SR_TYPE to backend's expected coordinate system
  IMAGE_FILE="test-deer-sr$SR_TYPE.png"
  curl -k "https://localhost:8443/arcgis/rest/services/Features/Environmental_admin/MapServer/export?dpi=96&transparent=true&format=png32&bbox=$DEER_BBOX&bboxSR=$SR_TYPE&imageSR=3424&size=512,512&f=image&layers=show:17" \
       -o "$IMAGE_FILE"
fi

echo "Test image saved to: $IMAGE_FILE"

# Capture all proxy logs
echo "Capturing proxy logs..."
LOG_FILE="test-proxy-$MAP_TYPE-sr$SR_TYPE.log"
podman logs wms-proxy-container > "$LOG_FILE"

echo "Logs saved to: $LOG_FILE"

# Stop the proxy
echo "Stopping proxy..."
make stop

echo "=== Test Complete ==="
echo "Files created:"
echo "  - $IMAGE_FILE (map image covering entire New Jersey using EPSG:$SR_TYPE coordinates)"
echo "  - $LOG_FILE (proxy logs)"
echo ""
echo "Coverage area: Entire New Jersey state"
echo "Coordinate transformation: EPSG:$SR_TYPE -> Backend's coordinate system (auto-detected)"
echo "Check the logs to see the REST URL that was generated and coordinate transformation!"
