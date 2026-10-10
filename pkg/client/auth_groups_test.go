package client

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAuthGroupsCRUD(t *testing.T) {
	var stored AuthGroupResponse

	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("Authorization = %q", r.Header.Get("Authorization"))
		}

		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/groups":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]AuthGroupResponse{stored})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/groups":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			var req AuthGroupCreateRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Fatalf("unmarshal create: %v", err)
			}
			mfaRequired := req.MfaRequired != nil && *req.MfaRequired
			stored = AuthGroupResponse{
				ID:          10,
				Name:        req.Name,
				IsSystem:    false,
				MfaRequired: mfaRequired,
				CreatedAt:   "2026-01-01T00:00:00.000Z",
				UpdatedAt:   "2026-01-01T00:00:00.000Z",
			}
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/v3/groups/"):
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodPatch && strings.HasPrefix(r.URL.Path, "/api/v3/groups/"):
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			var req AuthGroupUpdateRequest
			if err := json.Unmarshal(body, &req); err != nil {
				t.Fatalf("unmarshal update: %v", err)
			}
			if req.Name != nil {
				stored.Name = *req.Name
			}
			if req.MfaRequired != nil {
				stored.MfaRequired = *req.MfaRequired
			}
			stored.UpdatedAt = "2026-01-02T00:00:00.000Z"
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/v3/groups/"):
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})

	mfaRequired := true
	created, err := clt.CreateAuthGroup(AuthGroupCreateRequest{
		Name:        "secops",
		MfaRequired: &mfaRequired,
	})
	if err != nil {
		t.Fatalf("CreateAuthGroup: %v", err)
	}
	if created.Name != "secops" || !created.MfaRequired {
		t.Fatalf("created = %+v", created)
	}

	groups, err := clt.ListAuthGroups()
	if err != nil {
		t.Fatalf("ListAuthGroups: %v", err)
	}
	if len(groups) != 1 || groups[0].Name != "secops" {
		t.Fatalf("groups = %+v", groups)
	}

	got, err := clt.GetAuthGroup("secops")
	if err != nil {
		t.Fatalf("GetAuthGroup: %v", err)
	}
	if got.Name != "secops" {
		t.Fatalf("got = %+v", got)
	}

	newName := "platform-ops"
	mfaOff := false
	updated, err := clt.UpdateAuthGroup("secops", AuthGroupUpdateRequest{
		Name:        &newName,
		MfaRequired: &mfaOff,
	})
	if err != nil {
		t.Fatalf("UpdateAuthGroup: %v", err)
	}
	if updated.Name != "platform-ops" || updated.MfaRequired {
		t.Fatalf("updated = %+v", updated)
	}

	if err := clt.DeleteAuthGroup("platform-ops"); err != nil {
		t.Fatalf("DeleteAuthGroup: %v", err)
	}
}

func TestGetAuthGroupPathEscape(t *testing.T) {
	var requestedPath string

	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		requestedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(AuthGroupResponse{Name: "platform ops"})
	})

	if _, err := clt.GetAuthGroup("platform ops"); err != nil {
		t.Fatalf("GetAuthGroup: %v", err)
	}
	if requestedPath != "/api/v3/groups/platform%20ops" {
		t.Fatalf("path = %q, want %q", requestedPath, "/api/v3/groups/platform%20ops")
	}
}
