package microservices

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	ws "github.com/gorilla/websocket"
)

func TestDialControlWSUsesBearerAndPath(t *testing.T) {
	upgrader := ws.Upgrader{}
	ackReceived := make(chan bool, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != URLGetControlWSV1 {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer ws-token" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("failed to upgrade ws: %v", err)
		}
		defer conn.Close()

		if err := conn.WriteMessage(ws.BinaryMessage, []byte{CODE_CONTROL_SIGNAL}); err != nil {
			t.Fatalf("failed to write control signal: %v", err)
		}
		_, payload, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read ack: %v", err)
		}
		ackReceived <- len(payload) > 0 && payload[0] == CODE_ACK
	}))
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
	if err := os.WriteFile(tokenPath, []byte("ws-token"), 0o600); err != nil {
		t.Fatalf("failed to write token: %v", err)
	}

	opts := applyClientOptions(defaultClientOptions(),
		WithHost(host),
		WithPort(port),
		WithTLS(false),
		WithTokenPath(tokenPath),
	)
	client := newEdgeletAPIWsClient(opts)
	conn, err := client.dialControlWS()
	if err != nil {
		t.Fatalf("unexpected dial error: %v", err)
	}
	defer conn.Close()

	signals := make(chan byte, 1)
	go func() {
		_ = client.listenControl(conn, signals)
	}()

	select {
	case sig := <-signals:
		if sig != CODE_CONTROL_SIGNAL {
			t.Fatalf("unexpected signal opcode: %v", sig)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for control signal")
	}

	select {
	case ok := <-ackReceived:
		if !ok {
			t.Fatal("expected control ACK")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for ACK receipt")
	}
}
