package apps

import (
	"bytes"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestMicroserviceServiceAccountRoundTrip(t *testing.T) {
	const input = `
name: wasm-ms
application: demo
agent:
  name: edge-agent
serviceAccount:
  roleRef:
    kind: Role
    name: microservice
`
	var ms Microservice
	if err := yaml.Unmarshal([]byte(input), &ms); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if ms.ServiceAccount == nil {
		t.Fatal("expected serviceAccount to be set")
	}
	if ms.ServiceAccount.RoleRef.Kind != "Role" || ms.ServiceAccount.RoleRef.Name != "microservice" {
		t.Fatalf("unexpected roleRef: %+v", ms.ServiceAccount.RoleRef)
	}

	out, err := yaml.Marshal(ms)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(out)
	for _, want := range []string{"serviceAccount:", "roleRef:", "kind: Role", "name: microservice"} {
		if !strings.Contains(got, want) {
			t.Errorf("marshaled YAML missing %q:\n%s", want, got)
		}
	}
}

func TestMicroserviceDeployYAMLPreservesServiceAccount(t *testing.T) {
	ms := Microservice{
		Name:        "wasm-ms",
		Application: "demo",
		Agent:       MicroserviceAgent{Name: "edge-agent"},
		ServiceAccount: &MicroserviceServiceAccountRef{
			RoleRef: RoleRef{Kind: "Role", Name: "microservice"},
		},
	}
	file := IofogHeader{
		APIVersion: DefaultAPIVersion,
		Kind:       MicroserviceKind,
		Metadata: HeaderMetadata{
			Name: "demo/wasm-ms",
		},
		Spec: ms,
	}
	yamlBytes, err := yaml.Marshal(file)
	if err != nil {
		t.Fatalf("marshal IofogHeader: %v", err)
	}
	if !strings.Contains(string(yamlBytes), "serviceAccount:") {
		t.Fatalf("deploy YAML missing serviceAccount:\n%s", yamlBytes)
	}

	var fragment microserviceDeployServiceAccountFragment
	if err := yaml.Unmarshal(yamlBytes, &fragment); err != nil {
		t.Fatalf("unmarshal serviceAccount fragment: %v", err)
	}
	if fragment.Spec.ServiceAccount == nil || fragment.Spec.ServiceAccount.RoleRef.Name != "microservice" {
		t.Fatalf("serviceAccount lost in deploy YAML: %+v", fragment.Spec.ServiceAccount)
	}

	// Ensure CreateMicroserviceFromYAML-style reader still contains the block.
	if !bytes.Contains(yamlBytes, []byte("name: microservice")) {
		t.Fatalf("deploy YAML missing roleRef name:\n%s", yamlBytes)
	}
}

type microserviceDeployServiceAccountFragment struct {
	Spec microserviceDeployServiceAccountSpec `yaml:"spec"`
}

type microserviceDeployServiceAccountSpec struct {
	ServiceAccount *MicroserviceServiceAccountRef `yaml:"serviceAccount"`
}
