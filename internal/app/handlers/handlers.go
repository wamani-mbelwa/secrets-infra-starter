package handlers

import (
    "crypto/tls"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/example/wli-mtls-lab/internal/infra/tlsutil"
)

var (
    handshakes = promauto.NewCounter(prometheus.CounterOpts{
        Name: "mtls_handshakes_total",
        Help: "Total successful peer identity validations",
    })
)

func MetricsHandler() http.Handler { return promhttp.Handler() }

func Healthz(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "ok")
}

func Identity(w http.ResponseWriter, r *http.Request, allowed string) {
    cs := r.TLS
    if cs == nil {
        http.Error(w, "no TLS", http.StatusBadRequest); return
    }
    if err := tlsutil.VerifyPeerSPIFFE(*cs, allowed); err != nil {
        http.Error(w, err.Error(), http.StatusForbidden); return
    }
    handshakes.Inc()
    // Report peer cert notBefore/notAfter for visibility
    pc := cs.PeerCertificates[0]
    resp := map[string]any{
        "peer_spiffe_ids": tlsutil.ExtractSPIFFEIDs(pc),
        "not_before": pc.NotBefore.Format(time.RFC3339),
        "not_after":  pc.NotAfter.Format(time.RFC3339),
    }
    _ = json.NewEncoder(w).Encode(resp)
}

func LogTLSState(prefix string, cs *tls.ConnectionState) {
    if cs == nil || len(cs.PeerCertificates) == 0 { return }
    ids := tlsutil.ExtractSPIFFEIDs(cs.PeerCertificates[0])
    log.Printf("%s: peer ids=%v, vers=0x%x, cipher=0x%x", prefix, ids, cs.Version, cs.CipherSuite)
}
