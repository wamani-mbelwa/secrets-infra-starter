package config

import (
    "log"
    "os"
)

type Config struct {
    ServiceName   string
    ListenAddr    string
    MetricsAddr   string
    TLSCertFile   string
    TLSKeyFile    string
    TLSCAFile     string
    PeerAllowedID string
    OrderSvcURL   string
}

func FromEnv() Config {
    c := Config{
        ServiceName:   getenv("SERVICE_NAME", "svc"),
        ListenAddr:    getenv("LISTEN_ADDR", ":8080"),
        MetricsAddr:   getenv("METRICS_ADDR", ":9100"),
        TLSCertFile:   getenv("TLS_CERT_FILE", "certs/svc.crt"),
        TLSKeyFile:    getenv("TLS_KEY_FILE", "certs/svc.key"),
        TLSCAFile:     getenv("TLS_CA_FILE", "certs/ca.crt"),
        PeerAllowedID: getenv("PEER_ALLOWED_ID", ""),
        OrderSvcURL:   getenv("ORDER_SVC", ""),
    }
    log.Printf("config: %+v", c)
    return c
}

func getenv(k, def string) string {
    v := os.Getenv(k)
    if v == "" {
        return def
    }
    return v
}
