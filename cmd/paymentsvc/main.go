package main

import (
    "log"
    "net/http"
    "os"

    "github.com/example/wli-mtls-lab/internal/app/handlers"
    "github.com/example/wli-mtls-lab/internal/infra/config"
    "github.com/example/wli-mtls-lab/internal/infra/tlsutil"
    cl "github.com/example/wli-mtls-lab/pkg/client"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    cfg := config.FromEnv()

    ca, err := os.ReadFile(cfg.TLSCAFile); if err != nil { log.Fatal(err) }
    crt, err := os.ReadFile(cfg.TLSCertFile); if err != nil { log.Fatal(err) }
    key, err := os.ReadFile(cfg.TLSKeyFile); if err != nil { log.Fatal(err) }

    tcfgServer, err := tlsutil.TLSConfigServer(ca, crt, key); if err != nil { log.Fatal(err) }
    tcfgClient, err := tlsutil.TLSConfigClient(ca, crt, key); if err != nil { log.Fatal(err) }

    httpClient := cl.New(tcfgClient)

    mux := http.NewServeMux()
    mux.HandleFunc("/healthz", handlers.Healthz)
    mux.Handle("/metrics", promhttp.Handler())
    mux.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
        // Call ordersvc over mTLS, then verify its identity on their side.
        url := cfg.OrderSvcURL + "/identity"
        code, body, err := httpClient.Get(url)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadGateway); return
        }
        w.WriteHeader(code)
        w.Write(body)
    })
    srv := &http.Server{ Addr: cfg.ListenAddr, Handler: mux, TLSConfig: tcfgServer }
    log.Printf("paymentsvc listening on %s; calling %s", cfg.ListenAddr, cfg.OrderSvcURL)
    log.Fatal(srv.ListenAndServeTLS("", ""))
}
