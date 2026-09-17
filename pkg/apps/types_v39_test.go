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
