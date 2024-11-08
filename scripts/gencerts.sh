#!/usr/bin/env bash
set -euo pipefail
mkdir -p certs
pushd certs >/dev/null

# Create OpenSSL config for SAN URI (SPIFFE IDs)
cat > openssl.cnf <<'EOF'
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext
prompt             = no

[ req_distinguished_name ]
C  = US
ST = MN
L  = MSP
O  = Example
OU = Security
CN = Example Local CA

[ v3_ca ]
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid:always,issuer
basicConstraints = critical, CA:true
keyUsage = digitalSignature, keyEncipherment, keyCertSign, cRLSign

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
URI.1 = spiffe://example.org/ca
EOF

# CA key and cert
openssl genrsa -out ca.key 4096
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt -config openssl.cnf -extensions v3_ca

gen_leaf () {
  local name="$1"
  local spiffe="$2"
  cat > ${name}.cnf <<EOF
[ req ]
default_bits       = 2048
distinguished_name = req_distinguished_name
req_extensions     = req_ext
prompt             = no

[ req_distinguished_name ]
C  = US
ST = MN
L  = MSP
O  = Example
OU = Services
CN = ${name}

[ req_ext ]
subjectAltName = URI:${spiffe}
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth, clientAuth
EOF

  openssl genrsa -out ${name}.key 2048
  openssl req -new -key ${name}.key -out ${name}.csr -config ${name}.cnf
  openssl x509 -req -in ${name}.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out ${name}.crt -days 365 -sha256 -extfile ${name}.cnf -extensions req_ext
}

gen_leaf ordersvc spiffe://example.org/ns/default/sa/ordersvc
gen_leaf paymentsvc spiffe://example.org/ns/default/sa/paymentsvc

popd >/dev/null
echo "Certificates generated in ./certs"
