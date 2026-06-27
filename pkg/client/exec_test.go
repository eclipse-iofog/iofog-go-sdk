package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	ws "github.com/gorilla/websocket"
	json "github.com/json-iterator/go"
	"github.com/vmihailenco/msgpack/v5"
)

func TestParseExecSessionInfoFromActivation(t *testing.T) {
	t.Parallel()

	infoJSON, err := json.Marshal(ExecSessionInfo{
		SessionID:        "session-from-json",
		MicroserviceUUID: "ms-from-json",
	})
	if err != nil {
		t.Fatalf("marshal activation json: %v", err)
	}

	msg := &execWSMessage{
		Type:             ExecMessageActivation,
		Data:             infoJSON,
		MicroserviceUUID: "ms-top-level",
		ExecID:           "session-top-level",
		SessionID:        "session-top-level",
		Timestamp:        time.Now().UnixMilli(),
	}

	info, err := parseExecSessionInfo(msg)
	if err != nil {
		t.Fatalf("parseExecSessionInfo: %v", err)
	}
	if info.SessionID != "session-from-json" {
		t.Fatalf("expected sessionId from JSON data, got %q", info.SessionID)
	}
	if info.MicroserviceUUID != "ms-from-json" {
		t.Fatalf("expected microserviceUuid from JSON data, got %q", info.MicroserviceUUID)
	}
}

func TestParseExecSessionInfoUsesTopLevelFields(t *testing.T) {
	t.Parallel()

	msg := &execWSMessage{
		Type:             ExecMessageActivation,
		MicroserviceUUID: "ms-uuid",
		ExecID:           "exec-id",
	}

	info, err := parseExecSessionInfo(msg)
	if err != nil {
		t.Fatalf("parseExecSessionInfo: %v", err)
	}
	if info.SessionID != "exec-id" {
		t.Fatalf("expected exec-id as sessionId, got %q", info.SessionID)
	}
	if info.MicroserviceUUID != "ms-uuid" {
		t.Fatalf("expected ms-uuid, got %q", info.MicroserviceUUID)
	}
}

func TestExecWSURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		baseURL  string
		path     string
		expected string
	}{
		{
			name:     "http to ws",
			baseURL:  "http://localhost:51121/api/v3",
			path:     "/microservices/exec/uuid-1",
			expected: "ws://localhost:51121/api/v3/microservices/exec/uuid-1",
		},
		{
			name:     "https to wss system exec",
			baseURL:  "https://controller.example/api/v3",
			path:     "/microservices/system/exec/uuid-2",
			expected: "wss://controller.example/api/v3/microservices/system/exec/uuid-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			u, err := url.Parse(tt.baseURL)
			if err != nil {
				t.Fatalf("parse base url: %v", err)
			}
			clt := &Client{baseURL: u}
			got, err := clt.execWSURL(strings.TrimPrefix(tt.path, "/"))
			if err != nil {
				t.Fatalf("execWSURL: %v", err)
			}
			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestMapExecCloseError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		code   int
		reason string
		target error
	}{
		{"quota", ws.ClosePolicyViolation, "Maximum of 3 concurrent exec sessions allowed", ErrExecSessionQuotaExceeded},
		{"agent timeout", ws.ClosePolicyViolation, "Timeout waiting for agent connection", ErrExecAgentTimeout},
		{"not running", ws.ClosePolicyViolation, "Microservice is not running", ErrMicroserviceNotRunning},
		{"router unavailable", 1013, "Try again later", ErrExecRouterUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := mapExecCloseError(&ws.CloseError{Code: tt.code, Text: tt.reason})
			if err == nil {
				t.Fatal("expected mapped error")
			}
			if !strings.Contains(err.Error(), tt.target.Error()) {
				t.Fatalf("expected error containing %q, got %v", tt.target.Error(), err)
			}
		})
	}
}

func TestDialMicroserviceExec(t *testing.T) {
	const (
		msUUID    = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
		sessionID = "f9e8d7c6-b5a4-3210-fedc-ba9876543210"
	)

	upgrader := ws.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if r.URL.Path != "/api/v3/microservices/exec/"+msUUID {
			http.NotFound(w, r)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		activation, err := msgpack.Marshal(&execWSMessage{
			Type:             ExecMessageActivation,
			MicroserviceUUID: msUUID,
			ExecID:           sessionID,
			SessionID:        sessionID,
			Data: mustJSON(t, ExecSessionInfo{
				SessionID:        sessionID,
				MicroserviceUUID: msUUID,
			}),
		})
		if err != nil {
			return
		}
		if err := conn.WriteMessage(ws.BinaryMessage, activation); err != nil {
			return
		}

		stderr, err := msgpack.Marshal(&execWSMessage{
			Type:             ExecMessageStderr,
			MicroserviceUUID: msUUID,
			ExecID:           sessionID,
			SessionID:        sessionID,
			Data:             []byte("Agent connected. Interactive exec is ready."),
		})
		if err != nil {
			return
		}
		_ = conn.WriteMessage(ws.BinaryMessage, stderr)
		<-time.After(100 * time.Millisecond)
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

	session, err := clt.DialMicroserviceExec(msUUID)
	if err != nil {
		t.Fatalf("DialMicroserviceExec: %v", err)
	}
	defer session.Close()

	if session.SessionID != sessionID {
		t.Fatalf("expected sessionId %q, got %q", sessionID, session.SessionID)
	}
	if session.MicroserviceUUID != msUUID {
		t.Fatalf("expected microserviceUuid %q, got %q", msUUID, session.MicroserviceUUID)
	}

	if err := session.WriteStdin([]byte("echo hi\n")); err != nil {
		t.Fatalf("WriteStdin: %v", err)
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	return data
}
