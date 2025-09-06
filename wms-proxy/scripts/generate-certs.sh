#!/bin/bash

# Generate self-signed SSL certificates for the WMS proxy
# This script creates certificates suitable for development and testing

set -e

CERT_DIR="${CERT_DIR:-./certs}"
CERT_FILE="${CERT_FILE:-server.crt}"
KEY_FILE="${KEY_FILE:-server.key}"
DAYS="${DAYS:-365}"
COUNTRY="${COUNTRY:-US}"
STATE="${STATE:-State}"
CITY="${CITY:-City}"
ORG="${ORG:-WMS Proxy}"
OU="${OU:-IT Department}"
CN="${CN:-localhost}"

echo "Generating SSL certificates for WMS Proxy..."
echo "Certificate directory: $CERT_DIR"
echo "Certificate file: $CERT_FILE"
echo "Key file: $KEY_FILE"
echo "Valid for: $DAYS days"
echo "Common Name: $CN"

# Create certificate directory
mkdir -p "$CERT_DIR"

# Generate private key
echo "Generating private key..."
openssl genrsa -out "$CERT_DIR/$KEY_FILE" 2048

# Generate certificate signing request
echo "Generating certificate signing request..."
openssl req -new -key "$CERT_DIR/$KEY_FILE" -out "$CERT_DIR/server.csr" -subj "/C=$COUNTRY/ST=$STATE/L=$CITY/O=$ORG/OU=$OU/CN=$CN"

# Generate self-signed certificate with SAN (Subject Alternative Names)
echo "Generating self-signed certificate..."
cat > "$CERT_DIR/server.conf" << EOF
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = $COUNTRY
ST = $STATE
L = $CITY
O = $ORG
OU = $OU
CN = $CN

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = *.localhost
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

openssl x509 -req -in "$CERT_DIR/server.csr" -signkey "$CERT_DIR/$KEY_FILE" -out "$CERT_DIR/$CERT_FILE" -days "$DAYS" -extensions v3_req -extfile "$CERT_DIR/server.conf"

# Clean up temporary files
rm -f "$CERT_DIR/server.csr" "$CERT_DIR/server.conf"

# Set appropriate permissions (make both readable for container)
chmod 644 "$CERT_DIR/$KEY_FILE"
chmod 644 "$CERT_DIR/$CERT_FILE"

echo ""
echo "SSL certificates generated successfully!"
echo "Certificate: $CERT_DIR/$CERT_FILE"
echo "Private Key: $CERT_DIR/$KEY_FILE"
echo ""
echo "To use HTTPS with the proxy:"
echo "  export ENABLE_HTTPS=true"
echo "  export CERT_FILE=$PWD/$CERT_DIR/$CERT_FILE"
echo "  export KEY_FILE=$PWD/$CERT_DIR/$KEY_FILE"
echo ""
echo "Note: This is a self-signed certificate suitable for development only."
echo "For production, use certificates from a trusted Certificate Authority."
