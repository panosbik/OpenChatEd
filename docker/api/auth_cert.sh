#!/bin/sh

CERT_DIR="./auth.cert"
PRIV_KEY_FILE="$CERT_DIR/id_rsa"
PUB_KEY_FILE="$CERT_DIR/id_rsa.pub"

# Define the function to check if a file exists
file_exists() {
    if [ -f "$1" ]; then
        return 0 # The file exists
    else
        return 1 # The file does not exist
    fi
}

# Check if the cert directory already exists
if [ -d "$CERT_DIR" ]; then
    echo "cert directory already exists!"
else
    echo "Creating cert directory..."
    mkdir "$CERT_DIR"
fi

# Check if the private key and public key files already exist
if file_exists "$PRIV_KEY_FILE" && file_exists "$PUB_KEY_FILE"; then
    echo "Certificate files already exist!"
else
    echo "Generating certificate files..."
    openssl genrsa -out "$PRIV_KEY_FILE" 4096
    openssl rsa -in "$PRIV_KEY_FILE" -pubout -out "$PUB_KEY_FILE"
fi

echo "Certificate creation completed!"
