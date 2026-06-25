package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

const sampleAgentWithPlatformStatus = `{
	"uuid": "fdb5fca7-4cb1-4e4b-919a-43b116bc1f9d",
	"name": "edge-1",
	"routerMode": "interior",
	"natsMode": "server",
	"platformStatus": {
		"observedGeneration": 2,
		"phase": "Ready",
		"lastError": null,
		"lastTransitionAt": "2026-06-24T22:52:59.824Z",
		"conditions": [
			{"type": "RouterReady", "status": "True", "reason": "ReconcileComplete"},
			{"type": "NatsReady", "status": "True", "reason": "ReconcileComplete"}
		],
		"generation": 2
	}
}`

const sampleServiceReady = `{
	"name": "alert-dashboard",
	"type": "microservice",
	"resource": "10674117-42f3-473e-8e50-aa700d55ae46",
	"defaultBridge": "default-router",
	"bridgePort": 10025,
	"targetPort": 8082,
	"servicePort": 10025,
	"k8sType": null,
	"serviceEndpoint": "192.168.105.3",
	"tags": ["pot-edge-patterns"],
	"provisioningStatus": "ready",
	"provisioningError": null
}`

func TestAgentInfoPlatformStatusUnmarshal(t *testing.T) {
	var agent AgentInfo
	if err := json.Unmarshal([]byte(sampleAgentWithPlatformStatus), &agent); err != nil {
		t.Fatalf("Unmarshal AgentInfo: %v", err)
	}
	if agent.PlatformStatus == nil {
		t.Fatal("expected platformStatus")
	}
	if agent.PlatformStatus.Phase != PlatformReady {
		t.Fatalf("phase = %q, want %q", agent.PlatformStatus.Phase, PlatformReady)
	}
	if agent.PlatformStatus.Generation != 2 || agent.PlatformStatus.ObservedGeneration != 2 {
		t.Fatalf("generation mismatch: %+v", agent.PlatformStatus)
	}
	if agent.PlatformStatus.LastError != nil {
		t.Fatalf("lastError = %v, want nil", *agent.PlatformStatus.LastError)
	}
	if agent.PlatformStatus.LastTransitionAt == nil {
		t.Fatal("expected lastTransitionAt")
	}
	if len(agent.PlatformStatus.Conditions) != 2 {
		t.Fatalf("conditions len = %d, want 2", len(agent.PlatformStatus.Conditions))
	}
}

func TestServiceInfoProvisioningUnmarshal(t *testing.T) {
	var service ServiceInfo
	if err := json.Unmarshal([]byte(sampleServiceReady), &service); err != nil {
		t.Fatalf("Unmarshal ServiceInfo: %v", err)
	}
	if service.ProvisioningStatus != ProvisioningReady {
		t.Fatalf("provisioningStatus = %q, want %q", service.ProvisioningStatus, ProvisioningReady)
	}
	if service.ProvisioningError != nil {
		t.Fatalf("provisioningError = %v, want nil", *service.ProvisioningError)
	}
	if service.K8sType != nil {
		t.Fatalf("k8sType = %v, want nil", *service.K8sType)
	}
}

func testClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	baseURL, err := url.Parse(server.URL + "/api/v3")
	if err != nil {
		t.Fatalf("parse base URL: %v", err)
	}
	clt := &Client{
		baseURL:     baseURL,
		accessToken: "test-token",
		timeout:     5,
	}
	return clt
}

func TestReconcileAgent(t *testing.T) {
	const uuid = "fdb5fca7-4cb1-4e4b-919a-43b116bc1f9d"
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/iofog/"+uuid+"/reconcile" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Fatalf("authorization = %q", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"uuid":"` + uuid + `"}`))
	})

	resp, err := clt.ReconcileAgent(uuid)
	if err != nil {
		t.Fatalf("ReconcileAgent: %v", err)
	}
	if resp.UUID != uuid {
		t.Fatalf("UUID = %q, want %q", resp.UUID, uuid)
	}
}

func TestReconcileService(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/services/alert-dashboard/reconcile" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleServiceReady))
	})

	service, err := clt.ReconcileService("alert-dashboard")
	if err != nil {
		t.Fatalf("ReconcileService: %v", err)
	}
	if service.Name != "alert-dashboard" {
		t.Fatalf("name = %q", service.Name)
	}
	if service.ProvisioningStatus != ProvisioningReady {
		t.Fatalf("provisioningStatus = %q", service.ProvisioningStatus)
	}
}

func TestWaitForAgentPlatformReady(t *testing.T) {
	const uuid = "fdb5fca7-4cb1-4e4b-919a-43b116bc1f9d"
	calls := 0
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/iofog/"+uuid {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls < 2 {
			_, _ = w.Write([]byte(`{"uuid":"` + uuid + `","platformStatus":{"phase":"Progressing","generation":1,"observedGeneration":0}}`))
			return
		}
		_, _ = w.Write([]byte(sampleAgentWithPlatformStatus))
	})

	if err := clt.WaitForAgentPlatformReady(uuid, 5*time.Second); err != nil {
		t.Fatalf("WaitForAgentPlatformReady: %v", err)
	}
	if calls < 2 {
		t.Fatalf("expected at least 2 polls, got %d", calls)
	}
}

func TestWaitForServiceProvisioningReady(t *testing.T) {
	calls := 0
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/services/alert-dashboard" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls < 2 {
			_, _ = w.Write([]byte(`{"name":"alert-dashboard","provisioningStatus":"pending"}`))
			return
		}
		_, _ = w.Write([]byte(sampleServiceReady))
	})

	if err := clt.WaitForServiceProvisioningReady("alert-dashboard", 5*time.Second); err != nil {
		t.Fatalf("WaitForServiceProvisioningReady: %v", err)
	}
	if calls < 2 {
		t.Fatalf("expected at least 2 polls, got %d", calls)
	}
}

func TestWaitForAgentPlatformReadyFailed(t *testing.T) {
	const uuid = "fdb5fca7-4cb1-4e4b-919a-43b116bc1f9d"
	errMsg := "router cert failed"
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"uuid":"` + uuid + `","platformStatus":{"phase":"Failed","generation":1,"observedGeneration":1,"lastError":"` + errMsg + `"}}`))
	})

	err := clt.WaitForAgentPlatformReady(uuid, time.Second)
	if err == nil {
		t.Fatal("expected error for Failed phase")
	}
	if err.Error() == "" {
		t.Fatal("expected non-empty error message")
	}
}
