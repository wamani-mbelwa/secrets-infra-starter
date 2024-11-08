#!/usr/bin/env bash
set -euo pipefail
kind create cluster --name wli --wait 60s || true
echo "Kind cluster 'wli' is ready. Next: deploy SPIRE per deploy/spire/README.md then run 'make deploy'."
