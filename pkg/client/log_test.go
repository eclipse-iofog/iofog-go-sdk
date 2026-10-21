package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	ws "github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

func TestLogTailQuery(t *testing.T) {
	t.Parallel()

	got := logTailQuery(&LogTailOptions{
		Tail:   50,
		Follow: false,
		Since:  "2024-01-01T00:00:00Z",
	})
	if !strings.Contains(got, "tail=50") {
		t.Fatalf("expected tail=50 in %q", got)
	}
	if !strings.Contains(got, "follow=false") {
		t.Fatalf("expected follow=false in %q", got)
	}
	if !strings.Contains(got, "since=2024-01-01T00%3A00%3A00Z") {
		t.Fatalf("expected since query param in %q", got)
	}
}

func TestDialMicroserviceLogs(t *testing.T) {
	const msUUID = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

	upgrader := ws.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if r.URL.Path != "/api/v3/microservices/"+msUUID+"/logs" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("tail") != "10" {
			http.Error(w, "bad tail", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		start, err := msgpack.Marshal(&logWSMessage{
			Type:      LogMessageStart,
			ExecID:    "session-1",
			SessionID: "session-1",
		})
		if err != nil {
			return
		}
		if err := conn.WriteMessage(ws.BinaryMessage, start); err != nil {
			return
		}

		line, err := msgpack.Marshal(&logWSMessage{
			Type: LogMessageLine,
			Data: []byte("hello\n"),
		})
		if err != nil {
			return
		}
		_ = conn.WriteMessage(ws.BinaryMessage, line)

		stop, err := msgpack.Marshal(&logWSMessage{Type: LogMessageStop})
		if err != nil {
			return
		}
		_ = conn.WriteMessage(ws.BinaryMessage, stop)
	}))
	t.Cleanup(server.Close)

	baseURL, err := url.Parse(server.URL + "/api/v3")
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}

	clt := &Client{
		baseURL:     baseURL,
		accessToken: "test-token",
	}

	session, err := clt.DialMicroserviceLogs(msUUID, &LogTailOptions{Tail: 10, Follow: false})
	if err != nil {
		t.Fatalf("DialMicroserviceLogs: %v", err)
	}
	defer session.Close()

	frame, err := session.Read()
	if err != nil {
		t.Fatalf("Read start: %v", err)
	}
	if frame.Type != LogMessageStart {
		t.Fatalf("expected LOG_START, got type %d", frame.Type)
	}

	frame, err = session.Read()
	if err != nil {
		t.Fatalf("Read line: %v", err)
	}
	if frame.Type != LogMessageLine || string(frame.Data) != "hello\n" {
		t.Fatalf("unexpected log line frame: %+v", frame)
	}

	frame, err = session.Read()
	if err != nil {
		t.Fatalf("Read stop: %v", err)
	}
	if frame.Type != LogMessageStop {
		t.Fatalf("expected LOG_STOP, got type %d", frame.Type)
	}
}

func TestDialExecOnStatusLine(t *testing.T) {
	const (
		msUUID    = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
		sessionID = "f9e8d7c6-b5a4-3210-fedc-ba9876543210"
	)

	upgrader := ws.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		activation, _ := msgpack.Marshal(&execWSMessage{
			Type:             ExecMessageActivation,
			MicroserviceUUID: msUUID,
			SessionID:        sessionID,
			Data: mustJSON(t, ExecSessionInfo{
				SessionID:        sessionID,
				MicroserviceUUID: msUUID,
			}),
		})
		_ = conn.WriteMessage(ws.BinaryMessage, activation)

		waiting, _ := msgpack.Marshal(&execWSMessage{
			Type: ExecMessageStderr,
			Data: []byte("Waiting for agent connection..."),
		})
		_ = conn.WriteMessage(ws.BinaryMessage, waiting)

		ready, _ := msgpack.Marshal(&execWSMessage{
			Type: ExecMessageStderr,
			Data: []byte("Agent connected. Interactive exec is ready."),
		})
		_ = conn.WriteMessage(ws.BinaryMessage, ready)
	}))
	t.Cleanup(server.Close)

	baseURL, err := url.Parse(server.URL + "/api/v3")
	if err != nil {
		t.Fatalf("parse server url: %v", err)
	}

	clt := &Client{
		baseURL:     baseURL,
		accessToken: "test-token",
	}

	var lines []string
	session, err := clt.DialMicroserviceExecWithOptions(msUUID, &DialExecOptions{
		OnStatusLine: func(line string) {
			lines = append(lines, line)
		},
	})
	if err != nil {
		t.Fatalf("DialMicroserviceExecWithOptions: %v", err)
	}
	defer session.Close()

	if len(lines) != 2 {
		t.Fatalf("expected 2 status lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "Waiting for agent connection..." {
		t.Fatalf("unexpected first status line: %q", lines[0])
	}
}

func TestExecSessionCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	session := &ExecSession{}
	if err := session.Close(); err != nil {
		t.Fatalf("Close on nil conn: %v", err)
	}
	if err := session.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
}

func TestMapLogCloseError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		code   int
		reason string
		target error
	}{
		{"relay unavailable", 1013, "Relay unavailable for cross-replica session", ErrWsRelayUnavailable},
		{"legacy router reason", 1013, "Router unavailable for cross-replica session", ErrWsRelayUnavailable},
		{"server draining", ws.CloseGoingAway, "Server draining", ErrWsServerDraining},
		{"quota", ws.ClosePolicyViolation, "No available log session", ErrLogSessionUnavailable},
		{"agent timeout", ws.ClosePolicyViolation, "Timeout waiting for agent connection", ErrWsAgentTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := mapLogCloseError(&ws.CloseError{Code: tt.code, Text: tt.reason})
			if err == nil {
				t.Fatal("expected mapped error")
			}
			if !strings.Contains(err.Error(), tt.target.Error()) {
				t.Fatalf("expected error containing %q, got %v", tt.target.Error(), err)
			}
		})
	}
}

func TestLogSessionCloseIsIdempotent(t *testing.T) {
	t.Parallel()

	session := &LogSession{}
	if err := session.Close(); err != nil {
		t.Fatalf("Close on nil conn: %v", err)
	}
	if err := session.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
}
