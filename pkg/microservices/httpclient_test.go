package microservices

import (
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestHTTPClientGetConfigV3(t *testing.T) {
	caPEM, certPEM, keyPEM := generateCAAndServerCert(t)
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("failed to parse keypair: %v", err)
	}

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != URLGetConfigV1 {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"config": map[string]interface{}{
					"name": "svc",
				},
			},
		})
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
	server.StartTLS()
	defer server.Close()

	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse server URL: %v", err)
	}
	host, portStr, err := net.SplitHostPort(parsed.Host)
	if err != nil {
		t.Fatalf("failed to split host/port: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	caPath := filepath.Join(dir, "ca.crt")
	if err := os.WriteFile(tokenPath, []byte("test-token"), 0o600); err != nil {
		t.Fatalf("failed to write token: %v", err)
	}
	if err := os.WriteFile(caPath, caPEM, 0o600); err != nil {
		t.Fatalf("failed to write CA: %v", err)
	}

	opts := applyClientOptions(defaultClientOptions(),
		WithHost(host),
		WithPort(port),
		WithTLS(true),
		WithTokenPath(tokenPath),
		WithCAPath(caPath),
	)
	client := newIoFogHttpClient(opts)
	cfg, err := client.getConfig()
	if err != nil {
		t.Fatalf("unexpected getConfig error: %v", err)
	}
	if cfg["name"] != "svc" {
		t.Fatalf("unexpected config payload: %v", cfg)
	}
}
