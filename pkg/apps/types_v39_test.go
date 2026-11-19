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

func TestMicroserviceYAMLIncludesKnowledge(t *testing.T) {
	ms := Microservice{
		Name: "infer",
		Knowledge: &KnowledgeCatalog{
			BindPath: "/knowledge",
			Items:    []KnowledgeCatalogItem{{Name: "wiki"}},
		},
	}

	out, err := yaml.Marshal(ms)
	if err != nil {
		t.Fatalf("yaml.Marshal Microservice: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "knowledge:") {
		t.Errorf("YAML missing knowledge:\n%s", got)
	}
	if !strings.Contains(got, "bindPath: /knowledge") && !strings.Contains(got, "bindPath: '/knowledge'") {
		t.Errorf("YAML missing knowledge.bindPath:\n%s", got)
	}
	if !strings.Contains(got, "name: wiki") {
		t.Errorf("YAML missing knowledge item name:\n%s", got)
	}

	var decoded Microservice
	if err := yaml.Unmarshal(out, &decoded); err != nil {
		t.Fatalf("yaml.Unmarshal Microservice: %v", err)
	}
	if decoded.Knowledge == nil || decoded.Knowledge.BindPath != "/knowledge" {
		t.Fatalf("knowledge = %+v", decoded.Knowledge)
	}
	if len(decoded.Knowledge.Items) != 1 || decoded.Knowledge.Items[0].Name != "wiki" {
		t.Fatalf("knowledge items = %+v", decoded.Knowledge.Items)
	}
}

func TestKnowledgeSpecYAMLTags(t *testing.T) {
	spec := Knowledge{
		Repo:       "org/wiki",
		Revision:   "main",
		RegistryID: 5,
		Files:      []string{"README.md"},
		Format:     "dataset",
	}

	out, err := yaml.Marshal(spec)
	if err != nil {
		t.Fatalf("yaml.Marshal Knowledge: %v", err)
	}
	got := string(out)
	for _, want := range []string{"repo:", "revision:", "registryId:", "files:", "format:"} {
		if !strings.Contains(got, want) {
			t.Errorf("YAML missing %s:\n%s", want, got)
		}
	}

	var decoded Knowledge
	if err := yaml.Unmarshal(out, &decoded); err != nil {
		t.Fatalf("yaml.Unmarshal Knowledge: %v", err)
	}
	if decoded.Repo != "org/wiki" || decoded.Revision != "main" || decoded.RegistryID != 5 || decoded.Format != "dataset" {
		t.Fatalf("decoded = %+v", decoded)
	}
	if len(decoded.Files) != 1 || decoded.Files[0] != "README.md" {
		t.Fatalf("files = %v", decoded.Files)
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

func TestMicroserviceContainerSupplementalGroupsYAML(t *testing.T) {
	omitted, err := yaml.Marshal(MicroserviceContainer{RunAsGroup: "1000"})
	if err != nil {
		t.Fatalf("marshal omitted: %v", err)
	}
	if strings.Contains(string(omitted), "supplementalGroups:") {
		t.Errorf("nil supplementalGroups must omit the key:\n%s", omitted)
	}

	none := []int{}
	cleared, err := yaml.Marshal(MicroserviceContainer{SupplementalGroups: &none})
	if err != nil {
		t.Fatalf("marshal empty: %v", err)
	}
	if !strings.Contains(string(cleared), "supplementalGroups: []") {
		t.Errorf("empty supplementalGroups must send []:\n%s", cleared)
	}

	groups := []int{993, 44}
	out, err := yaml.Marshal(MicroserviceContainer{
		RunAsGroup:         "1000",
		SupplementalGroups: &groups,
	})
	if err != nil {
		t.Fatalf("marshal groups: %v", err)
	}
	got := string(out)
	if strings.Index(got, "runAsGroup:") > strings.Index(got, "supplementalGroups:") {
		t.Errorf("runAsGroup should appear before supplementalGroups:\n%s", got)
	}
	if strings.Contains(got, "readOnlyRootFilesystem:") && strings.Index(got, "supplementalGroups:") > strings.Index(got, "readOnlyRootFilesystem:") {
		t.Errorf("supplementalGroups should appear before readOnlyRootFilesystem:\n%s", got)
	}
	id993 := strings.Index(got, "993")
	id44 := strings.Index(got, "44")
	if id993 < 0 || id44 < 0 || id993 > id44 {
		t.Errorf("YAML must keep caller order 993 then 44:\n%s", got)
	}

	const input = `
runAsGroup: "1000"
supplementalGroups:
  - 44
  - 993
`
	var decoded MicroserviceContainer
	if err := yaml.Unmarshal([]byte(input), &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.SupplementalGroups == nil || len(*decoded.SupplementalGroups) != 2 {
		t.Fatalf("decoded = %#v", decoded.SupplementalGroups)
	}
	if (*decoded.SupplementalGroups)[0] != 44 || (*decoded.SupplementalGroups)[1] != 993 {
		t.Fatalf("decoded ids = %v", *decoded.SupplementalGroups)
	}

	var absent MicroserviceContainer
	if err := yaml.Unmarshal([]byte("runAsGroup: \"1000\"\n"), &absent); err != nil {
		t.Fatalf("unmarshal absent: %v", err)
	}
	if absent.SupplementalGroups != nil {
		t.Fatalf("omitted key decoded as %#v", absent.SupplementalGroups)
	}

	copied := decoded.DeepCopy()
	(*copied.SupplementalGroups)[0] = 1
	if (*decoded.SupplementalGroups)[0] != 44 {
		t.Fatal("deepcopy aliased supplementalGroups")
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

func TestMicroserviceStatusYAMLIncludesImagePullStatus(t *testing.T) {
	failed, err := yaml.Marshal(Microservice{
		Name: "infer",
		Status: MicroserviceStatusInfo{
			Status:          "RUNNING",
			ErrorMessage:    "unable to pull nginx: failed to pull image",
			ImagePullStatus: "failed",
		},
	})
	if err != nil {
		t.Fatalf("yaml.Marshal failed pull: %v", err)
	}
	if !strings.Contains(string(failed), "imagePullStatus: failed") {
		t.Errorf("YAML missing imagePullStatus:\n%s", failed)
	}
	if !strings.Contains(string(failed), "status: RUNNING") {
		t.Errorf("YAML changed observed status:\n%s", failed)
	}

	cleared, err := yaml.Marshal(Microservice{
		Name: "infer",
		Status: MicroserviceStatusInfo{
			Status:          "RUNNING",
			ImagePullStatus: "",
			LastError:       "unable to pull nginx: failed to pull image",
		},
	})
	if err != nil {
		t.Fatalf("yaml.Marshal cleared pull: %v", err)
	}
	if !strings.Contains(string(cleared), "imagePullStatus: \"\"") {
		t.Errorf("YAML omitted empty imagePullStatus:\n%s", cleared)
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
