# Cert Rotation Runbook
- Confirm current leaf cert TTL via /identity response (not_after).
- If TTL < 6h, trigger rotation workflow (SPIRE default rotation or step-ca renew).
- Ensure servers reload keys via zero-downtime (liveness/readiness probes).
- Verify using smoke: paymentsvc -> ordersvc /identity shows new not_after.
