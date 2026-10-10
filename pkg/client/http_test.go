package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func startTLSTestServer(t *testing.T) (*httptest.Server, []byte) {
	t.Helper()

	caPEM, certPEM, keyPEM := generateCAAndServerCert(t)
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("failed to parse keypair: %v", err)
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
			"versions": map[string]string{
				"controller": "3.8.0",
			},
		})
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
	server.StartTLS()
	t.Cleanup(server.Close)
	return server, caPEM
}

func TestHTTPDoDefaultTLSAllowsSelfSigned(t *testing.T) {
	server, _ := startTLSTestServer(t)

	hd := httpDo{timeout: 5}
	_, err := hd.do("GET", server.URL, map[string]string{"Content-Type": "application/json"}, nil)
	if err != nil {
		t.Fatalf("expected default TLS to accept self-signed cert: %v", err)
	}
}

func TestHTTPDoTLSConfigWithRootCAsVerifiesCert(t *testing.T) {
	server, caPEM := startTLSTestServer(t)

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		t.Fatal("failed to append CA cert")
	}

	hd := httpDo{
		timeout:   5,
		tlsConfig: &tls.Config{RootCAs: pool},
	}
	_, err := hd.do("GET", server.URL, map[string]string{"Content-Type": "application/json"}, nil)
	if err != nil {
		t.Fatalf("expected verified TLS request to succeed: %v", err)
	}
}

func TestHTTPDoTLSConfigInsecureSkipVerify(t *testing.T) {
	server, _ := startTLSTestServer(t)

	hd := httpDo{
		timeout:   5,
		tlsConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402 -- explicit insecure mode
	}
	_, err := hd.do("GET", server.URL, map[string]string{"Content-Type": "application/json"}, nil)
	if err != nil {
		t.Fatalf("expected insecure TLS request to succeed: %v", err)
	}
}

func TestHTTPDoTLSConfigWithoutRootCAsRejectsSelfSigned(t *testing.T) {
	server, _ := startTLSTestServer(t)

	hd := httpDo{
		timeout:   5,
		tlsConfig: &tls.Config{},
	}
	_, err := hd.do("GET", server.URL, map[string]string{"Content-Type": "application/json"}, nil)
	if err == nil {
		t.Fatal("expected TLS verification failure for self-signed cert")
	}
	if !strings.Contains(err.Error(), "certificate") {
		t.Fatalf("expected certificate verification error, got: %v", err)
	}
}
