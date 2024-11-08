#!/usr/bin/env bash
set -euo pipefail
curl -sk https://localhost:8082/healthz || true
curl -sk https://localhost:8081/healthz || true
echo "Attempting mTLS call paymentsvc -> ordersvc (requires compose stack)."
curl -sk https://localhost:8082/pay || true
