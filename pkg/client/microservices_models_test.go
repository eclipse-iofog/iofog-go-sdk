package client

import (
	"net/http"
	"strings"
	"testing"
)

const testMicroserviceUUID = "10674117-42f3-473e-8e50-aa700d55ae46"

func TestPatchMicroserviceModels(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v3/microservices/"+testMicroserviceUUID+"/models" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		req := readJSONObject(t, r)
		if req["bindPath"] != "/models" {
			t.Fatalf("bindPath = %v", req["bindPath"])
		}
		if req["permissions"] != "ro" {
			t.Fatalf("permissions = %v", req["permissions"])
		}
		items, ok := req["items"].([]any)
		if !ok || len(items) != 1 {
			t.Fatalf("items = %v", req["items"])
		}
		item, ok := items[0].(map[string]any)
		if !ok {
			t.Fatalf("item type = %T", items[0])
		}
		if item["name"] != "llama" {
			t.Fatalf("item name = %v", item["name"])
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := clt.PatchMicroserviceModels(testMicroserviceUUID, MicroserviceCatalog{
		BindPath:    "/models",
		Permissions: "ro",
		Items:       []MicroserviceCatalogItem{{Name: "llama"}},
	})
	if err != nil {
		t.Fatalf("PatchMicroserviceModels: %v", err)
	}
}

func TestCreateMicroserviceFromYAMLField(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/microservices/yaml":
			if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
				t.Fatal("CreateMicroserviceFromYAML must not send JSON")
			}
			assertMultipartFormFile(t, r, "microservice")
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"uuid":"` + testMicroserviceUUID + `"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/microservices/"+testMicroserviceUUID:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"uuid":"` + testMicroserviceUUID + `","name":"infer"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})

	info, err := clt.CreateMicroserviceFromYAML(strings.NewReader("kind: Microservice\nmetadata:\n  name: infer\n"))
	if err != nil {
		t.Fatalf("CreateMicroserviceFromYAML: %v", err)
	}
	if info.UUID != testMicroserviceUUID {
		t.Fatalf("uuid = %q", info.UUID)
	}
}
