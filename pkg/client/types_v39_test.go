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

func TestMicroserviceInfoSupplementalGroups(t *testing.T) {
	raw := []byte(`{
		"uuid":"ms-1",
		"name":"infer",
		"runAsGroup":"1000",
		"supplementalGroups":[44,993]
	}`)

	var info MicroserviceInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("Unmarshal MicroserviceInfo: %v", err)
	}
	if info.RunAsGroup != "1000" {
		t.Fatalf("runAsGroup = %q", info.RunAsGroup)
	}
	if len(info.SupplementalGroups) != 2 || info.SupplementalGroups[0] != 44 || info.SupplementalGroups[1] != 993 {
		t.Fatalf("supplementalGroups = %v", info.SupplementalGroups)
	}

	dup := []byte(`{"uuid":"ms-1","name":"infer","supplementalGroups":[44,44]}`)
	var fromDup MicroserviceInfo
	if err := json.Unmarshal(dup, &fromDup); err != nil {
		t.Fatalf("Unmarshal duplicates: %v", err)
	}
	if len(fromDup.SupplementalGroups) != 2 {
		t.Fatalf("duplicates collapsed: %v", fromDup.SupplementalGroups)
	}

	emptyRaw := []byte(`{"uuid":"ms-1","name":"infer","supplementalGroups":[]}`)
	var empty MicroserviceInfo
	if err := json.Unmarshal(emptyRaw, &empty); err != nil {
		t.Fatalf("Unmarshal empty: %v", err)
	}
	if empty.SupplementalGroups == nil || len(empty.SupplementalGroups) != 0 {
		t.Fatalf("empty list = %#v", empty.SupplementalGroups)
	}
	body, err := json.Marshal(empty)
	if err != nil {
		t.Fatalf("Marshal empty: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("Unmarshal marshaled empty: %v", err)
	}
	groups, ok := got["supplementalGroups"].([]any)
	if !ok || len(groups) != 0 {
		t.Fatalf("marshaled supplementalGroups = %s", body)
	}

	fraction := []byte(`{"uuid":"ms-1","name":"infer","supplementalGroups":[44.5]}`)
	if err := json.Unmarshal(fraction, &MicroserviceInfo{}); err == nil {
		t.Fatal("44.5 must not unmarshal into supplementalGroups")
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

func TestMicroserviceStatusInfoUnmarshalImagePullStatus(t *testing.T) {
	failed := []byte(`{
		"uuid":"ms-1",
		"name":"infer",
		"status":{
			"status":"RUNNING",
			"containerId":"d8ed513372a8",
			"errorMessage":"unable to pull docker.io/library/nginx:edgelet-pull-miss-absent: failed to pull image",
			"imagePullStatus":"failed",
			"lastError":"unable to pull docker.io/library/nginx:edgelet-pull-miss-absent: failed to pull image",
			"lastErrorAt":1730000000000,
			"restartCount":0
		}
	}`)

	var running MicroserviceInfo
	if err := json.Unmarshal(failed, &running); err != nil {
		t.Fatalf("Unmarshal running pull miss: %v", err)
	}
	if running.Status.Status != "RUNNING" {
		t.Fatalf("status.status = %q", running.Status.Status)
	}
	if running.Status.ImagePullStatus != "failed" {
		t.Fatalf("status.imagePullStatus = %q", running.Status.ImagePullStatus)
	}
	if running.Status.RestartCount != 0 {
		t.Fatalf("status.restartCount = %d", running.Status.RestartCount)
	}

	queued := []byte(`{
		"uuid":"ms-2",
		"name":"infer",
		"status":{"status":"QUEUED","imagePullStatus":"failed","errorMessage":"unable to pull nginx: failed to pull image"}
	}`)
	var neverStarted MicroserviceInfo
	if err := json.Unmarshal(queued, &neverStarted); err != nil {
		t.Fatalf("Unmarshal queued pull miss: %v", err)
	}
	if neverStarted.Status.Status != "QUEUED" || neverStarted.Status.ImagePullStatus != "failed" {
		t.Fatalf("status = %q imagePullStatus = %q", neverStarted.Status.Status, neverStarted.Status.ImagePullStatus)
	}
	if neverStarted.Status.ContainerID != "" {
		t.Fatalf("containerId = %q", neverStarted.Status.ContainerID)
	}

	cleared := []byte(`{"uuid":"ms-3","status":{"status":"RUNNING","imagePullStatus":"","errorMessage":"","lastError":"unable to pull nginx: failed to pull image"}}`)
	var recovered MicroserviceInfo
	if err := json.Unmarshal(cleared, &recovered); err != nil {
		t.Fatalf("Unmarshal cleared pull: %v", err)
	}
	if recovered.Status.ImagePullStatus != "" || recovered.Status.ErrorMessage != "" {
		t.Fatalf("imagePullStatus = %q errorMessage = %q", recovered.Status.ImagePullStatus, recovered.Status.ErrorMessage)
	}
	if recovered.Status.LastError == "" {
		t.Fatal("lastError was cleared")
	}

	omitted := []byte(`{"uuid":"ms-4","status":{"status":"RUNNING"}}`)
	var older MicroserviceInfo
	if err := json.Unmarshal(omitted, &older); err != nil {
		t.Fatalf("Unmarshal omitted imagePullStatus: %v", err)
	}
	if older.Status.ImagePullStatus != "" {
		t.Fatalf("omitted imagePullStatus = %q", older.Status.ImagePullStatus)
	}

	encoded, err := json.Marshal(older.Status)
	if err != nil {
		t.Fatalf("Marshal status: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(encoded, &raw); err != nil {
		t.Fatalf("Unmarshal encoded status: %v", err)
	}
	pull, ok := raw["imagePullStatus"]
	if !ok {
		t.Fatalf("encoded status omitted imagePullStatus: %s", encoded)
	}
	if pull != "" {
		t.Fatalf("encoded imagePullStatus = %v", pull)
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
				"imagePullStatus":"",
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
	if st.ImagePullStatus != "" {
		t.Fatalf("status.imagePullStatus = %q", st.ImagePullStatus)
	}
}

func TestAgentInfoUnmarshalHostStatus(t *testing.T) {
	raw := []byte(`{
		"uuid":"fog-1",
		"name":"edge-1",
		"systemCpus":8,
		"systemTotalMemory":17179869184,
		"systemAvailableMemory":8589934592,
		"systemTotalDisk":107374182400,
		"systemAvailableDisk":53687091200,
		"systemTotalCpu":12.5,
		"systemOs":"linux",
		"systemOsVersion":"Ubuntu 22.04",
		"systemKernelVersion":"6.8.0",
		"diskUsage":1.25
	}`)

	var agent AgentInfo
	if err := json.Unmarshal(raw, &agent); err != nil {
		t.Fatalf("Unmarshal AgentInfo: %v", err)
	}
	if agent.SystemCpus != 8 {
		t.Fatalf("systemCpus = %d", agent.SystemCpus)
	}
	if agent.SystemTotalMemory != 17179869184 {
		t.Fatalf("systemTotalMemory = %d", agent.SystemTotalMemory)
	}
	if agent.SystemAvailableMemory != 8589934592 {
		t.Fatalf("systemAvailableMemory = %d", agent.SystemAvailableMemory)
	}
	if agent.SystemTotalDisk != 107374182400 {
		t.Fatalf("systemTotalDisk = %d", agent.SystemTotalDisk)
	}
	if agent.SystemAvailableDisk != 53687091200 {
		t.Fatalf("systemAvailableDisk = %d", agent.SystemAvailableDisk)
	}
	if agent.SystemTotalCPU != 12.5 {
		t.Fatalf("systemTotalCpu = %v", agent.SystemTotalCPU)
	}
	if agent.SystemOs != "linux" || agent.SystemOsVersion != "Ubuntu 22.04" || agent.SystemKernelVersion != "6.8.0" {
		t.Fatalf("system os fields = %q %q %q", agent.SystemOs, agent.SystemOsVersion, agent.SystemKernelVersion)
	}
	if agent.DiskUsage != 1.25 {
		t.Fatalf("diskUsage = %v", agent.DiskUsage)
	}
}

func TestAgentInfoUnmarshalV39FieldsWithoutHAL(t *testing.T) {
	// modelLastUpdate and knowledgeLastUpdate are Unix milliseconds.
	raw := []byte(`{
		"uuid":"fog-1",
		"name":"edge-1",
		"runtimeClasses":"[{\"name\":\"nvidia\"}]",
		"activeModels":2,
		"modelLastUpdate":1710000000000,
		"knowledgeStatus":"[{\"name\":\"wiki\"}]",
		"activeKnowledge":3,
		"knowledgeLastUpdate":1710000000456
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
	if agent.ModelLastUpdate != 1710000000000 {
		t.Fatalf("modelLastUpdate = %d", agent.ModelLastUpdate)
	}
	if agent.KnowledgeStatus != `[{"name":"wiki"}]` {
		t.Fatalf("knowledgeStatus = %q", agent.KnowledgeStatus)
	}
	if agent.ActiveKnowledge != 3 {
		t.Fatalf("activeKnowledge = %d", agent.ActiveKnowledge)
	}
	if agent.KnowledgeLastUpdate != 1710000000456 {
		t.Fatalf("knowledgeLastUpdate = %d", agent.KnowledgeLastUpdate)
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
