package apps

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func assertRegistryYAMLRoundTrip(t *testing.T, name, input string, wantContains string, wantNotContains []string) {
	t.Helper()
	var ms Microservice
	if err := yaml.Unmarshal([]byte(input), &ms); err != nil {
		t.Fatalf("%s: unmarshal: %v", name, err)
	}
	out, err := yaml.Marshal(ms)
	if err != nil {
		t.Fatalf("%s: marshal: %v", name, err)
	}
	got := string(out)
	if !strings.Contains(got, wantContains) {
		t.Errorf("%s: marshaled YAML missing %q:\n%s", name, wantContains, got)
	}
	for _, bad := range wantNotContains {
		if strings.Contains(got, bad) {
			t.Errorf("%s: marshaled YAML must not contain %q:\n%s", name, bad, got)
		}
	}
}

func TestMicroserviceImagesRegistryIntRoundTrip(t *testing.T) {
	const input = `
name: clickhouse
images:
  registry: 5
  amd64: dhi.io/clickhouse-server:26.7
`
	assertRegistryYAMLRoundTrip(t, "registry int", input, "registry: 5", []string{`registry: "5"`})
	if ms := unmarshalMicroservice(t, input); ms.Images.Registry.Int() != 5 {
		t.Fatalf("expected registry id 5, got %d", ms.Images.Registry.Int())
	}
}

func TestMicroserviceImagesRegistryAliasRoundTrip(t *testing.T) {
	const input = `
name: app
images:
  registry: remote
  amd64: ghcr.io/example/app:latest
`
	assertRegistryYAMLRoundTrip(t, "registry alias", input, "registry: 1", []string{"registry: remote", `registry: "1"`})
}

func TestMicroserviceImagesRegistryNumericStringRoundTrip(t *testing.T) {
	const input = `
name: app
images:
  registry: "10"
  amd64: ghcr.io/example/app:latest
`
	assertRegistryYAMLRoundTrip(t, "registry numeric string", input, "registry: 10", []string{`registry: "10"`})
}

func TestMicroserviceImagesRegistryInvalid(t *testing.T) {
	const input = `
name: app
images:
  registry: oci
  amd64: ghcr.io/example/app:latest
`
	var ms Microservice
	if err := yaml.Unmarshal([]byte(input), &ms); err == nil {
		t.Fatal("expected error for invalid registry alias oci")
	}
}

func TestCatalogItemRegistryRoundTrip(t *testing.T) {
	const input = `
id: 1
registry: 5
name: demo
`
	var item CatalogItem
	if err := yaml.Unmarshal([]byte(input), &item); err != nil {
		t.Fatalf("unmarshal CatalogItem: %v", err)
	}
	if item.Registry.Int() != 5 {
		t.Fatalf("expected registry id 5, got %d", item.Registry.Int())
	}
	out, err := yaml.Marshal(item)
	if err != nil {
		t.Fatalf("marshal CatalogItem: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "registry: 5") {
		t.Errorf("marshaled YAML missing registry: 5:\n%s", got)
	}
	if strings.Contains(got, `registry: "5"`) {
		t.Errorf("marshaled YAML must not quote registry:\n%s", got)
	}
}

func unmarshalMicroservice(t *testing.T, input string) Microservice {
	t.Helper()
	var ms Microservice
	if err := yaml.Unmarshal([]byte(input), &ms); err != nil {
		t.Fatalf("unmarshal Microservice: %v", err)
	}
	return ms
}
