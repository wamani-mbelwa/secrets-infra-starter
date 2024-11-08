package client

import (
    "crypto/tls"
    "io"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    mtlsRequests = promauto.NewCounter(prometheus.CounterOpts{
        Name: "mtls_client_requests_total",
        Help: "Total mTLS client requests",
    })
)

type MTLSClient struct {
    hc *http.Client
}

func New(tcfg *tls.Config) *MTLSClient {
    return &MTLSClient{
        hc: &http.Client{
            Transport: &http.Transport{ TLSClientConfig: tcfg },
            Timeout: 5 * time.Second,
        },
    }
}

func (c *MTLSClient) Get(url string) (int, []byte, error) {
    mtlsRequests.Inc()
    resp, err := c.hc.Get(url)
    if err != nil { return 0, nil, err }
    defer resp.Body.Close()
    b, _ := io.ReadAll(resp.Body)
    return resp.StatusCode, b, nil
}
