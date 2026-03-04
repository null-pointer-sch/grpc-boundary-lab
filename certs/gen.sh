#!/usr/bin/env bash
set -e

# Switch to script directory
cd "$(dirname "$0")"

echo "Generating Certificates..."

# 1. Generate Local Root CA
echo "[1] Generating Root CA (ca.key, ca.crt)..."
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 -out ca.crt \
    -subj "/C=US/ST=State/L=City/O=GrpcBoundaryLab/OU=CA/CN=Local Root CA"

generate_leaf() {
    local name=$1
    local san1=$2
    local san2=$3

    echo "[$name] Generating private key & CSR..."
    openssl genrsa -out "${name}.key" 2048
    
    # Create CSR
    openssl req -new -key "${name}.key" -out "${name}.csr" \
        -subj "/C=US/ST=State/L=City/O=GrpcBoundaryLab/OU=Service/CN=${san1}"

    # Create v3 ext config for SANs
    cat > "${name}_ext.cnf" <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = ${san1}
DNS.2 = ${san2}
IP.1 = 127.0.0.1
EOF

    echo "[$name] Signing certificate with Root CA..."
    openssl x509 -req -in "${name}.csr" -CA ca.crt -CAkey ca.key -CAcreateserial \
        -out "${name}.crt" -days 365 -sha256 -extfile "${name}_ext.cnf"

    # Cleanup intermediate files
    rm "${name}.csr" "${name}_ext.cnf"
}

# 2. Generate Frontend Cert
generate_leaf "frontend" "localhost" "frontend"

# 3. Generate Gateway Cert
generate_leaf "gateway" "localhost" "gateway"

# 4. Generate Backend Cert
generate_leaf "backend" "localhost" "backend"

echo "Certificates generated successfully!"
ls -la
