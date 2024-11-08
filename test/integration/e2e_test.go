//go:build integration

package integration

import "testing"

func TestE2E(t *testing.T) {
    t.Skip("Requires local env (docker compose or kind + SPIRE). See README.")
}
