package client

import (
	"encoding/json"
	"testing"
)

func TestRegistryCreateRequestMarshalHF(t *testing.T) {
	body, err := json.Marshal(RegistryCreateRequest{
		URL:      "https://huggingface.co",
		Type:     "hf",
		CA:       "base64pem",
		Insecure: true,
	})
	if err != nil {
		t.Fatalf("Marshal RegistryCreateRequest: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got["type"] != "hf" {
		t.Fatalf("type = %v, want hf", got["type"])
	}
	if got["ca"] != "base64pem" {
		t.Fatalf("ca = %v", got["ca"])
	}
	insecure, ok := got["insecure"].(bool)
	if !ok || !insecure {
		t.Fatalf("insecure = %v, want true", got["insecure"])
	}
	if _, ok := got["isSecure"]; ok {
		t.Fatalf("isSecure must not be marshaled: %s", body)
	}
	if _, ok := got["requiresCert"]; ok {
		t.Fatalf("requiresCert must not be marshaled: %s", body)
	}
}

func TestMicroserviceInfoUnmarshalV39Fields(t *testing.T) {
	raw := []byte(`{
		"uuid":"ms-1",
		"name":"infer",
		"runAsGroup":"1000",
		"entrypoint":["/usr/bin/python"],
		"cpus":1.5,
		"models":{"bindPath":"/models","permissions":"ro","items":[{"name":"llama"}]},
		"status":{"podId":"pod-abc","status":"RUNNING"},
		"tmpfs":[{"containerPath":"/tmp","size":64}],
		"devices":[{"hostPath":"/dev/nvidia0","containerPath":"/dev/nvidia0","permissions":"rwm"}]
	}`)

	var info MicroserviceInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("Unmarshal MicroserviceInfo: %v", err)
	}
	if info.RunAsGroup != "1000" {
		t.Fatalf("runAsGroup = %q", info.RunAsGroup)
	}
	if len(info.Entrypoint) != 1 || info.Entrypoint[0] != "/usr/bin/python" {
		t.Fatalf("entrypoint = %v", info.Entrypoint)
	}
	if info.Models == nil || info.Models.BindPath != "/models" {
		t.Fatalf("models = %+v", info.Models)
	}
	if info.Status.PodID != "pod-abc" {
		t.Fatalf("status.podId = %q", info.Status.PodID)
	}
	if info.Cpus != 1.5 {
		t.Fatalf("cpus = %v", info.Cpus)
	}
	if len(info.Tmpfs) != 1 || info.Tmpfs[0].ContainerPath != "/tmp" || info.Tmpfs[0].Size != 64 {
		t.Fatalf("tmpfs = %+v", info.Tmpfs)
	}
	if len(info.Devices) != 1 || info.Devices[0].HostPath != "/dev/nvidia0" {
		t.Fatalf("devices = %+v", info.Devices)
	}
}

func TestAgentInfoUnmarshalV39FieldsWithoutHAL(t *testing.T) {
	raw := []byte(`{
		"uuid":"fog-1",
		"name":"edge-1",
		"runtimeClasses":"[{\"name\":\"nvidia\"}]",
		"activeModels":2,
		"modelLastUpdate":1710000000
	}`)

	var agent AgentInfo
	if err := json.Unmarshal(raw, &agent); err != nil {
		t.Fatalf("Unmarshal AgentInfo: %v", err)
	}
	if agent.RuntimeClasses != `[{"name":"nvidia"}]` {
		t.Fatalf("runtimeClasses = %q", agent.RuntimeClasses)
	}
	if agent.ActiveModels != 2 {
		t.Fatalf("activeModels = %d", agent.ActiveModels)
	}
	if agent.ModelLastUpdate != 1710000000 {
		t.Fatalf("modelLastUpdate = %d", agent.ModelLastUpdate)
	}
}

func TestFogLinkRequestMarshalFogUUIDsOnly(t *testing.T) {
	body, err := json.Marshal(FogLinkRequest{FogUUIDs: []string{testFogUUID}})
	if err != nil {
		t.Fatalf("Marshal FogLinkRequest: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("keys = %v, want fogUuids only", got)
	}
	if _, ok := got["name"]; ok {
		t.Fatalf("name must not be marshaled: %s", body)
	}
	uuids, ok := got["fogUuids"].([]any)
	if !ok || len(uuids) != 1 || uuids[0] != testFogUUID {
		t.Fatalf("fogUuids = %v", got["fogUuids"])
	}
}
