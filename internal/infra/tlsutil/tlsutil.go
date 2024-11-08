package tlsutil

import (
    "crypto/tls"
    "crypto/x509"
    "errors"
    "fmt"
    "net/http"
    "strings"
)

// ExtractSPIFFEIDs returns all URI SANs with spiffe:// scheme from a cert.
func ExtractSPIFFEIDs(cert *x509.Certificate) []string {
    var out []string
    for _, uri := range cert.URIs {
        if uri.Scheme == "spiffe" {
            out = append(out, uri.String())
        }
    }
    return out
}

// VerifyPeerSPIFFE enforces that at least one peer certificate contains the allowed SPIFFE ID.
func VerifyPeerSPIFFE(cs tls.ConnectionState, allowed string) error {
    if allowed == "" {
        return errors.New("no allowed SPIFFE ID configured")
    }
    if len(cs.PeerCertificates) == 0 {
        return errors.New("no peer certificates")
    }
    for _, pc := range cs.PeerCertificates {
        ids := ExtractSPIFFEIDs(pc)
        for _, id := range ids {
            if strings.EqualFold(id, allowed) {
                return nil
            }
        }
    }
    return fmt.Errorf("peer SPIFFE ID not allowed; expected=%s", allowed)
}

func TLSConfigServer(caPEM, certPEM, keyPEM []byte) (*tls.Config, error) {
    cert, err := tls.X509KeyPair(certPEM, keyPEM)
    if err != nil { return nil, err }
    pool := x509.NewCertPool()
    if ok := pool.AppendCertsFromPEM(caPEM); !ok { return nil, errors.New("bad ca") }
    cfg := &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientCAs:    pool,
        ClientAuth:   tls.RequireAndVerifyClientCert,
        MinVersion:   tls.VersionTLS13,
    }
    return cfg, nil
}

func TLSConfigClient(caPEM, certPEM, keyPEM []byte) (*tls.Config, error) {
    cert, err := tls.X509KeyPair(certPEM, keyPEM)
    if err != nil { return nil, err }
    pool := x509.NewCertPool()
    if ok := pool.AppendCertsFromPEM(caPEM); !ok { return nil, errors.New("bad ca") }
    cfg := &tls.Config{
        Certificates:       []tls.Certificate{cert},
        RootCAs:            pool,
        MinVersion:         tls.VersionTLS13,
        Renegotiation:      tls.RenegotiateNever,
    }
    return cfg, nil
}

func NewHTTPClient(tcfg *tls.Config) *http.Client {
    tr := &http.Transport{ TLSClientConfig: tcfg }
    return &http.Client{ Transport: tr }
}
