package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

const testAgentUUID = "fdb5fca7-4cb1-4e4b-919a-43b116bc1f9d"

func readRequestBody(t *testing.T, r *http.Request) string {
	t.Helper()
	if r.Body == nil {
		return ""
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(body)
}

func assertVersionCommandRequest(t *testing.T, r *http.Request, command, wantBody string) {
	t.Helper()
	if r.Method != http.MethodPost {
		t.Fatalf("method = %q, want POST", r.Method)
	}
	wantPath := "/api/v3/iofog/" + testAgentUUID + "/version/" + command
	if r.URL.Path != wantPath {
		t.Fatalf("path = %q, want %q", r.URL.Path, wantPath)
	}
	if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
		t.Fatalf("authorization = %q", auth)
	}
	if got := readRequestBody(t, r); got != wantBody {
		t.Fatalf("body = %q, want %q", got, wantBody)
	}
}

func TestUpgradeNodeWithoutSemver(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertVersionCommandRequest(t, r, "upgrade", "")
		w.WriteHeader(http.StatusNoContent)
	})

	if err := clt.UpgradeNode(testAgentUUID, nil); err != nil {
		t.Fatalf("UpgradeNode: %v", err)
	}
}

func TestUpgradeNodeWithSemver(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertVersionCommandRequest(t, r, "upgrade", `{"semver":"3.2.0"}`)
		w.WriteHeader(http.StatusNoContent)
	})

	semver := "3.2.0"
	if err := clt.UpgradeNode(testAgentUUID, &semver); err != nil {
		t.Fatalf("UpgradeNode: %v", err)
	}
}

func TestRollbackNodeWithSemver(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertVersionCommandRequest(t, r, "rollback", `{"semver":"3.1.0"}`)
		w.WriteHeader(http.StatusNoContent)
	})

	semver := "3.1.0"
	if err := clt.RollbackNode(testAgentUUID, &semver); err != nil {
		t.Fatalf("RollbackNode: %v", err)
	}
}

func TestSetNodeVersionCommandWithoutSemver(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertVersionCommandRequest(t, r, "upgrade", "")
		w.WriteHeader(http.StatusNoContent)
	})

	if err := clt.SetNodeVersionCommand(testAgentUUID, "upgrade", nil); err != nil {
		t.Fatalf("SetNodeVersionCommand: %v", err)
	}
}

func TestSetNodeVersionCommandWithSemver(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertVersionCommandRequest(t, r, "rollback", `{"semver":"v1.0.0-beta.2"}`)
		w.WriteHeader(http.StatusNoContent)
	})

	semver := "v1.0.0-beta.2"
	if err := clt.SetNodeVersionCommand(testAgentUUID, "rollback", &SetNodeVersionCommandRequest{
		Semver: &semver,
	}); err != nil {
		t.Fatalf("SetNodeVersionCommand: %v", err)
	}
}

func TestSetNodeVersionCommandInvalidSemver(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message": "ValidationError",
		})
	})

	semver := "not-a-semver"
	err := clt.SetNodeVersionCommand(testAgentUUID, "upgrade", &SetNodeVersionCommandRequest{
		Semver: &semver,
	})
	if err == nil {
		t.Fatal("expected error for invalid semver")
	}

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected HTTPError, got %T: %v", err, err)
	}
	if httpErr.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want %d", httpErr.Code, http.StatusBadRequest)
	}
}

func TestUpgradeNodeNotReady(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"INVALID_VERSION_COMMAND_UPGRADE"}`))
	})

	err := clt.UpgradeNode(testAgentUUID, nil)
	if err == nil {
		t.Fatal("expected error when node is not ready to upgrade")
	}

	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected HTTPError, got %T: %v", err, err)
	}
	if httpErr.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want %d", httpErr.Code, http.StatusBadRequest)
	}
}
