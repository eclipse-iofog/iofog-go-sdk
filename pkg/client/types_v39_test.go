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
	if info.Status.LastError != "" || info.Status.LastErrorAt != 0 || info.Status.RestartCount != 0 {
		t.Fatalf("status extras must be zero when omitted: %+v", info.Status)
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

func TestMicroserviceVolumeMappingInfoRoundTripScope(t *testing.T) {
	raw := []byte(`{
		"uuid":"ms-1",
		"name":"infer",
		"volumeMappings":[{
			"hostDestination":"nodered-config",
			"containerDestination":"/data",
			"accessMode":"rw",
			"type":"volume",
			"scope":"shared"
		}]
	}`)

	var info MicroserviceInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("Unmarshal MicroserviceInfo: %v", err)
	}
	if len(info.Volumes) != 1 {
		t.Fatalf("volumeMappings len = %d", len(info.Volumes))
	}
	vol := info.Volumes[0]
	if vol.HostDestination != "nodered-config" || vol.Type != "volume" || vol.Scope != "shared" {
		t.Fatalf("volumeMapping = %+v", vol)
	}

	body, err := json.Marshal(vol)
	if err != nil {
		t.Fatalf("Marshal MicroserviceVolumeMappingInfo: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got["scope"] != "shared" {
		t.Fatalf("scope = %v, want shared", got["scope"])
	}
	if _, ok := got["id"]; ok {
		t.Fatalf("id must not be marshaled: %s", body)
	}

	empty, err := json.Marshal(MicroserviceVolumeMappingInfo{
		HostDestination:      "data",
		ContainerDestination: "/var/lib/app",
		AccessMode:           "rw",
		Type:                 "volume",
	})
	if err != nil {
		t.Fatalf("Marshal empty scope: %v", err)
	}
	var omitted map[string]any
	if err := json.Unmarshal(empty, &omitted); err != nil {
		t.Fatalf("Unmarshal empty scope: %v", err)
	}
	if _, ok := omitted["scope"]; ok {
		t.Fatalf("empty scope must be omitted: %s", empty)
	}
}

func TestMicroserviceStatusInfoUnmarshalErrorExtras(t *testing.T) {
	raw := []byte(`{
		"uuid":"ms-1",
		"name":"infer",
		"status":{
			"status":"RUNNING",
			"errorMessage":"",
			"lastError":"OOMKilled",
			"lastErrorAt":1710000000123,
			"restartCount":2
		}
	}`)

	var info MicroserviceInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("Unmarshal MicroserviceInfo: %v", err)
	}
	if info.Status.LastError != "OOMKilled" {
		t.Fatalf("status.lastError = %q", info.Status.LastError)
	}
	if info.Status.LastErrorAt != 1710000000123 {
		t.Fatalf("status.lastErrorAt = %d", info.Status.LastErrorAt)
	}
	if info.Status.RestartCount != 2 {
		t.Fatalf("status.restartCount = %d", info.Status.RestartCount)
	}
}

func TestApplicationInfoUnmarshalMicroserviceStatusExtras(t *testing.T) {
	raw := []byte(`{
		"name":"demo",
		"microservices":[{
			"uuid":"ms-1",
			"name":"infer",
			"status":{
				"status":"FAILED",
				"errorMessage":"crash loop",
				"lastError":"crash loop",
				"lastErrorAt":1710000000456,
				"restartCount":1
			}
		}]
	}`)

	var app ApplicationInfo
	if err := json.Unmarshal(raw, &app); err != nil {
		t.Fatalf("Unmarshal ApplicationInfo: %v", err)
	}
	if len(app.Microservices) != 1 {
		t.Fatalf("microservices len = %d", len(app.Microservices))
	}
	st := app.Microservices[0].Status
	if st.ErrorMessage != "crash loop" {
		t.Fatalf("status.errorMessage = %q", st.ErrorMessage)
	}
	if st.LastError != "crash loop" {
		t.Fatalf("status.lastError = %q", st.LastError)
	}
	if st.LastErrorAt != 1710000000456 {
		t.Fatalf("status.lastErrorAt = %d", st.LastErrorAt)
	}
	if st.RestartCount != 1 {
		t.Fatalf("status.restartCount = %d", st.RestartCount)
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
