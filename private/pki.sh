#!/bin/bash

DOMAIN="api.gpsd.com"
TOP_DIR="private"
CA_DIR="ca"
CERT_DIR="certs"
KEY_SIZE=2048
VALIDITY_DAYS=365

mkdir -p $TOP_DIR/$CA_DIR $TOP_DIR/$CERT_DIR

echo "Creating Certificate Authority (CA)..."

openssl genrsa -out $TOP_DIR/$CA_DIR/ca.key $KEY_SIZE

openssl req -x509 -new -nodes -key $TOP_DIR/$CA_DIR/ca.key -sha256 -days $VALIDITY_DAYS \
    -subj "/CN=GPSD_CA" -out $TOP_DIR/$CA_DIR/ca.crt

echo "CA created: $TOP_DIR/$CA_DIR/ca.crt"

echo "Generating certificate for domain: $DOMAIN..."

openssl genrsa -out $TOP_DIR/$CERT_DIR/$DOMAIN.key $KEY_SIZE

openssl req -new -key $TOP_DIR/$CERT_DIR/$DOMAIN.key -out $TOP_DIR/$CERT_DIR/$DOMAIN.csr \
    -subj "/CN=$DOMAIN"

cat > $TOP_DIR/$CERT_DIR/$DOMAIN.ext <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = $DOMAIN
EOF

openssl x509 -req -in $TOP_DIR/$CERT_DIR/$DOMAIN.csr -CA $TOP_DIR/$CA_DIR/ca.crt -CAkey $TOP_DIR/$CA_DIR/ca.key \
    -CAcreateserial -out $TOP_DIR/$CERT_DIR/$DOMAIN.crt -days $VALIDITY_DAYS -sha256 \
    -extfile $TOP_DIR/$CERT_DIR/$DOMAIN.ext

echo "Certificate issued: $TOP_DIR/$CERT_DIR/$DOMAIN.crt"

rm -f $TOP_DIR/$CERT_DIR/$DOMAIN.csr $TOP_DIR/$CERT_DIR/$DOMAIN.ext

echo "Certificates are ready:"
echo "CA Certificate: $TOP_DIR/$CA_DIR/ca.crt"
echo "Server Certificate: $TOP_DIR/$CERT_DIR/$DOMAIN.crt"
echo "Server Private Key: $TOP_DIR/$CERT_DIR/$DOMAIN.key"

openssl x509 -in $TOP_DIR/$CERT_DIR/$DOMAIN.crt -text -noout