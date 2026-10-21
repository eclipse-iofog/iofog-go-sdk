package client

import (
	"encoding/json"
	"testing"
)

func TestConfigMapInfoUnmarshal(t *testing.T) {
	raw := `{
		"id": 1,
		"name": "my-config",
		"data": {"key1": "value1"},
		"immutable": false,
		"useVault": true,
		"created_at": "2026-07-09T00:00:00.000Z",
		"updated_at": "2026-07-09T00:00:00.000Z"
	}`

	var info ConfigMapInfo
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		t.Fatalf("Unmarshal ConfigMapInfo: %v", err)
	}
	if info.Name != "my-config" {
		t.Fatalf("Name = %q, want my-config", info.Name)
	}
	if !info.UseVault {
		t.Fatal("UseVault = false, want true")
	}
	if info.CreatedAt != "2026-07-09T00:00:00.000Z" {
		t.Fatalf("CreatedAt = %q, want 2026-07-09T00:00:00.000Z", info.CreatedAt)
	}
}

func TestConfigMapCreateRequestMarshal(t *testing.T) {
	useVault := false
	immutable := false
	req := ConfigMapCreateRequest{
		Name:      "my-config",
		Data:      map[string]string{"key1": "value1"},
		UseVault:  &useVault,
		Immutable: &immutable,
	}

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal ConfigMapCreateRequest: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("Unmarshal request body: %v", err)
	}
	gotUseVault, ok := got["useVault"].(bool)
	if !ok || gotUseVault {
		t.Fatalf("useVault = %v, want false", got["useVault"])
	}
	immutableVal, ok := got["immutable"].(bool)
	if !ok || immutableVal {
		t.Fatalf("immutable = %v, want false", got["immutable"])
	}
}

func TestConfigMapUpdateRequestMarshalImmutableFalse(t *testing.T) {
	immutable := false
	req := ConfigMapUpdateRequest{
		Name:      "test-configmap",
		Data:      map[string]string{"key": "value"},
		Immutable: &immutable,
	}

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal ConfigMapUpdateRequest: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("Unmarshal request body: %v", err)
	}
	gotImmutable, ok := got["immutable"].(bool)
	if !ok || gotImmutable {
		t.Fatalf("immutable = %v, want false", got["immutable"])
	}
}

func TestConfigMapUpdateRequestOmitUseVault(t *testing.T) {
	req := ConfigMapUpdateRequest{
		Data: map[string]string{"key1": "new-value"},
	}

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal ConfigMapUpdateRequest: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("Unmarshal request body: %v", err)
	}
	if _, ok := got["useVault"]; ok {
		t.Fatalf("useVault should be omitted on patch, got %v", got["useVault"])
	}
}
