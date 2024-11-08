# Workload Identity & mTLS Lab (Go + SPIFFE/SPIRE)

Production-ready starter implementing workload identity and mTLS between services using SPIFFE/SPIRE principles.
Includes: Clean Architecture layout, Docker Compose path with local CA, kind/K8s path with SPIRE (manifests),
tiny Go client lib for identity-aware requests, metrics/logging, tests, and CI.

## TL;DR (Docker Compose, local CA)
```bash
# 1) Generate local CA and leaf certs with SPIFFE URI SANs
./scripts/gencerts.sh

# 2) Start services with mTLS enforced
docker compose up --build -d

# 3) Health checks
curl -sk https://localhost:8081/healthz
curl -sk https://localhost:8082/healthz

# 4) paymentsvc -> ordersvc mTLS call (server verifies peer SPIFFE ID)
curl -sk https://localhost:8082/pay
```

## K8s Path (kind + SPIRE)
- Bring up cluster: `make kind-up`
- Install SPIRE (example manifests under `deploy/spire/`; see README there).
- Apply kustomize overlays: `make deploy`
- Run E2E tests: `make e2e`

## Architecture
- Two services: `ordersvc` (server) and `paymentsvc` (client + server).
- mTLS required; server validates peer via SPIFFE ID allowlist (URI SAN check).
- `pkg/client/` exposes a minimal Go client helper for mTLS + identity assertions.
- Observability: Prometheus metrics and structured logs; Grafana dashboards JSON included.

## Clean Architecture layout
- `internal/domain` – pure types
- `internal/app` – core use-cases/handlers
- `internal/adapters` – http/grpc wire, metrics/logging, tls verification
- `internal/infra` – config, cert store
- `pkg/client` – client lib used by paymentsvc (and others)

## Tests
- Unit tests for handlers and SPIFFE URI parsing / allowlist checks.

## Security & DX
- Short TTLs recommended in prod; rotation runbooks in `docs/runbooks/`.
- No secrets committed; see `.env.example` and `scripts/gencerts.sh` for dev.

## CI
- See `.github/workflows/ci.yml`: lint, build and unit tests
```
