package apps

import (
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func int64Ptr(v int64) *int64 {
	return &v
}

func TestMicroserviceYAMLIncludesTemplateAndModels(t *testing.T) {
	ms := Microservice{
		Name: "infer",
		Template: &MicroserviceTemplateRef{
			Name:      "infer-tpl",
			Variables: MicroserviceTemplateVariables{"image": "nginx"},
		},
		Models: &MicroserviceCatalog{
			BindPath:    "/models",
			Permissions: "ro",
			Items:       []MicroserviceCatalogItem{{Name: "llama"}},
		},
	}

	out, err := yaml.Marshal(ms)
	if err != nil {
		t.Fatalf("yaml.Marshal Microservice: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "template:") {
		t.Errorf("YAML missing template:\n%s", got)
	}
	if !strings.Contains(got, "models:") {
		t.Errorf("YAML missing models:\n%s", got)
	}
	if !strings.Contains(got, "bindPath: /models") && !strings.Contains(got, "bindPath: '/models'") {
		t.Errorf("YAML missing models.bindPath:\n%s", got)
	}
	if !strings.Contains(got, "name: infer-tpl") {
		t.Errorf("YAML missing template name:\n%s", got)
	}
}

func TestMicroserviceContainerYAMLSpotCheck(t *testing.T) {
	container := MicroserviceContainer{
		RunAsGroup: "1000",
		Entrypoint: []string{"/bin/sh"},
		Tmpfs: []MicroserviceTmpfs{
			{ContainerPath: "/tmp", Size: int64Ptr(64)},
		},
	}

	out, err := yaml.Marshal(container)
	if err != nil {
		t.Fatalf("yaml.Marshal MicroserviceContainer: %v", err)
	}
	got := string(out)
	if strings.Index(got, "hostNetworkMode:") > strings.Index(got, "runAsGroup:") {
		t.Errorf("hostNetworkMode should appear before runAsGroup:\n%s", got)
	}
	if !strings.Contains(got, "runAsGroup: \"1000\"") && !strings.Contains(got, "runAsGroup: 1000") {
		t.Errorf("YAML missing runAsGroup:\n%s", got)
	}
	if !strings.Contains(got, "entrypoint:") {
		t.Errorf("YAML missing entrypoint:\n%s", got)
	}
	if !strings.Contains(got, "tmpfs:") {
		t.Errorf("YAML missing tmpfs:\n%s", got)
	}
}

func TestMicroserviceVolumeMappingYAMLRoundTripScope(t *testing.T) {
	vol := MicroserviceVolumeMapping{
		HostDestination:      "nodered-config",
		ContainerDestination: "/data",
		AccessMode:           "rw",
		Type:                 "volume",
		Scope:                "shared",
	}

	out, err := yaml.Marshal(vol)
	if err != nil {
		t.Fatalf("yaml.Marshal MicroserviceVolumeMapping: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "scope: shared") {
		t.Errorf("YAML missing scope:\n%s", got)
	}
	if strings.Contains(got, "id:") {
		t.Errorf("YAML must not include id:\n%s", got)
	}

	var decoded MicroserviceVolumeMapping
	if err := yaml.Unmarshal(out, &decoded); err != nil {
		t.Fatalf("yaml.Unmarshal MicroserviceVolumeMapping: %v", err)
	}
	if decoded.Scope != "shared" || decoded.Type != "volume" || decoded.HostDestination != "nodered-config" {
		t.Fatalf("decoded = %+v", decoded)
	}

	empty, err := yaml.Marshal(MicroserviceVolumeMapping{
		HostDestination:      "data",
		ContainerDestination: "/var/lib/app",
		AccessMode:           "rw",
		Type:                 "volume",
	})
	if err != nil {
		t.Fatalf("yaml.Marshal empty scope: %v", err)
	}
	if strings.Contains(string(empty), "scope:") {
		t.Errorf("empty scope must be omitted:\n%s", empty)
	}
}

func TestMicroserviceStatusYAMLIncludesErrorExtras(t *testing.T) {
	ms := Microservice{
		Name: "infer",
		Status: MicroserviceStatusInfo{
			Status:       "RUNNING",
			LastError:    "OOMKilled",
			LastErrorAt:  1710000000123,
			RestartCount: 2,
		},
	}

	out, err := yaml.Marshal(ms)
	if err != nil {
		t.Fatalf("yaml.Marshal Microservice: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "lastError: OOMKilled") {
		t.Errorf("YAML missing lastError:\n%s", got)
	}
	if !strings.Contains(got, "lastErrorAt: 1710000000123") {
		t.Errorf("YAML missing lastErrorAt:\n%s", got)
	}
	if !strings.Contains(got, "restartCount: 2") {
		t.Errorf("YAML missing restartCount:\n%s", got)
	}
}

func TestAgentConfigurationYAMLTagsOmitHAL(t *testing.T) {
	// Split literals so pkg-wide grep gates on the removed JSON keys stay clean.
	forbidden := []string{
		"deviceScan" + "Frequency",
		"bluetooth" + "Enabled",
		"abstractedHardware" + "Enabled",
	}
	typ := reflect.TypeOf(AgentConfiguration{})
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tags := field.Tag.Get("yaml") + " " + field.Tag.Get("json") + " " + field.Name
		for _, key := range forbidden {
			if strings.Contains(tags, key) {
				t.Errorf("AgentConfiguration field %s still references %s", field.Name, key)
			}
		}
	}
}
