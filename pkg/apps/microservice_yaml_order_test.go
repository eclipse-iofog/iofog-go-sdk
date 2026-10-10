package apps

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func assertYAMLKeysInOrder(t *testing.T, yamlStr string, keys ...string) {
	t.Helper()
	last := -1
	for _, key := range keys {
		idx := strings.Index(yamlStr, key+":")
		if idx < 0 {
			t.Errorf("missing key %q in:\n%s", key, yamlStr)
			continue
		}
		if idx <= last {
			t.Errorf("key %q is out of order (idx=%d, prev=%d):\n%s", key, idx, last, yamlStr)
		}
		last = idx
	}
}

func TestMicroserviceCanonicalYAMLFieldOrder(t *testing.T) {
	const input = `
uuid: db811baf-0b8d-4c55-936e-bcb1a4ac2b46
application: test-app
name: clickhouse
agent:
  name: lima
images:
  registry: 5
  arm64: dhi.io/clickhouse-server:26.7
  amd64: dhi.io/clickhouse-server:26.7
natsConfig:
  natsAccess: false
models:
  bindPath: /models
  items:
    - name: llama
container:
  hostNetworkMode: false
  isPrivileged: false
  runAsUser: ""
  runAsGroup: ""
  readOnlyRootFilesystem: false
  ulimits:
    nofile:
      soft: 262144
      hard: 262144
  shmSize: 1024
  ports:
    - internal: 8123
      external: 8123
schedule: 50
config:
  myKey: value
serviceAccount:
  roleRef:
    kind: Role
    name: microservice
`
	var ms Microservice
	if err := yaml.Unmarshal([]byte(input), &ms); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	out, err := yaml.Marshal(ms)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(out)

	assertYAMLKeysInOrder(t, got,
		"uuid", "application", "name", "agent", "images",
		"natsConfig", "models", "container", "schedule", "config", "serviceAccount",
	)

	imagesBlock := got[strings.Index(got, "images:"):]
	assertYAMLKeysInOrder(t, imagesBlock, "registry", "arm64", "amd64")

	containerBlock := got[strings.Index(got, "container:"):]
	assertYAMLKeysInOrder(t, containerBlock,
		"hostNetworkMode", "isPrivileged", "readOnlyRootFilesystem",
		"ulimits", "shmSize", "ports",
	)
}

func TestMicroserviceDeployRoundTripRegistryAliasAndServiceAccount(t *testing.T) {
	const input = `
application: test-app
name: wasm-ms
agent:
  name: edge-agent
images:
  registry: remote
  amd64: ghcr.io/example/app:latest
serviceAccount:
  roleRef:
    kind: Role
    name: microservice
schedule: 50
config: {}
`
	var ms Microservice
	if err := yaml.Unmarshal([]byte(input), &ms); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	out, err := yaml.Marshal(ms)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(out)

	if !strings.Contains(got, "registry: 1") {
		t.Errorf("expected registry: 1 after alias round-trip:\n%s", got)
	}
	if strings.Contains(got, `registry: "1"`) {
		t.Errorf("registry must not be quoted:\n%s", got)
	}
	if !strings.Contains(got, "serviceAccount:") || !strings.Contains(got, "name: microservice") {
		t.Errorf("serviceAccount block missing after round-trip:\n%s", got)
	}

	assertYAMLKeysInOrder(t, got,
		"application", "name", "agent", "images", "schedule", "config", "serviceAccount",
	)
}

func TestMicroserviceImagesOmitsZeroCatalogID(t *testing.T) {
	images := MicroserviceImages{
		Registry: RegistryRef(5),
		ARM64:    "dhi.io/clickhouse-server:26.7",
		AMD64:    "dhi.io/clickhouse-server:26.7",
	}
	out, err := yaml.Marshal(images)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(out)
	if strings.Contains(got, "catalogId:") {
		t.Errorf("catalogId should be omitted when zero:\n%s", got)
	}
	assertYAMLKeysInOrder(t, got, "registry", "arm64", "amd64")
}

func TestMicroserviceContainerFieldOrderWhenSet(t *testing.T) {
	container := MicroserviceContainer{
		HostNetworkMode:        false,
		IsPrivileged:           false,
		RunAsUser:              "1000",
		RunAsGroup:             "1000",
		ReadOnlyRootFilesystem: false,
		Ports:                  []MicroservicePortMapping{{Internal: 80, External: 8080}},
		Commands:               []string{"/bin/app"},
	}
	out, err := yaml.Marshal(container)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(out)
	assertYAMLKeysInOrder(t, got,
		"hostNetworkMode", "isPrivileged", "runAsUser", "runAsGroup", "readOnlyRootFilesystem",
		"ports", "commands",
	)
}

func TestCatalogItemRegistryFirstOrder(t *testing.T) {
	item := CatalogItem{
		ID:       1,
		Registry: RegistryRef(5),
		ARM64:    "img:arm64",
		AMD64:    "img:amd64",
		Name:     "demo",
	}
	out, err := yaml.Marshal(item)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(out)
	assertYAMLKeysInOrder(t, got, "id", "registry", "arm64", "amd64", "name")
}
