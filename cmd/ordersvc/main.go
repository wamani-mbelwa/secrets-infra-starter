package main

import (
    "crypto/tls"
    "crypto/x509"
    "embed"
    "io/fs"
    "log"
    "net/http"
    "os"

    "github.com/example/wli-mtls-lab/internal/app/handlers"
    "github.com/example/wli-mtls-lab/internal/infra/config"
    "github.com/example/wli-mtls-lab/internal/infra/tlsutil"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed ../../certs/*
var certFS embed.FS

func loadPEM(path string) ([]byte, error) {
    b, err := os.ReadFile(path)
    if err == nil { return b, nil }
    // fallback to embedded certs if available
    bb, err2 := fs.ReadFile(certFS, path)
    if err2 == nil { return bb, nil }
    return nil, err
}

func main() {
    cfg := config.FromEnv()

    ca, _ := loadPEM("certs/ca.crt")
    crt, _ := loadPEM(os.Getenv("TLS_CERT_FILE"))
    key, _ := loadPEM(os.Getenv("TLS_KEY_FILE"))
    tcfg, err := tlsutil.TLSConfigServer(ca, crt, key)
    if err != nil { log.Fatal(err) }
    tcfg.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
        // Let standard verification happen via ClientAuth, then we add policy in handler
        return nil
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", handlers.Healthz)
    mux.Handle("/metrics", promhttp.Handler())
    mux.HandleFunc("/identity", func(w http.ResponseWriter, r *http.Request) {
        handlers.LogTLSState("ordersvc", r.TLS)
        handlers.Identity(w, r, cfg.PeerAllowedID)
    })
    srv := &http.Server{
        Addr:      cfg.ListenAddr,
        Handler:   mux,
        TLSConfig: tcfg,
    }
    log.Printf("ordersvc listening on %s", cfg.ListenAddr)
    log.Fatal(srv.ListenAndServeTLS("", ""))
}
