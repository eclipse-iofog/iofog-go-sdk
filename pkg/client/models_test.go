package client

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

const testModelName = "llama"
const testModelUUID = "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
const testFogUUID = "fdb5fca7-4cb1-4e4b-919a-43b116bc1f9d"
const testBlockingMSUUID = "10674117-42f3-473e-8e50-aa700d55ae46"

func sampleModelJSON() string {
	return `{
		"uuid":"` + testModelUUID + `",
		"name":"` + testModelName + `",
		"repo":"org/llama",
		"revision":"main",
		"registryId":5,
		"format":"gguf"
	}`
}

func TestListModels(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/models" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"models":[` + sampleModelJSON() + `]}`))
	})

	resp, err := clt.ListModels()
	if err != nil {
		t.Fatalf("ListModels: %v", err)
	}
	if len(resp.Models) != 1 {
		t.Fatalf("models len = %d, want 1", len(resp.Models))
	}
	if resp.Models[0].Name != testModelName || resp.Models[0].Repo != "org/llama" || resp.Models[0].RegistryID != 5 {
		t.Fatalf("unexpected model: %+v", resp.Models[0])
	}
}

func TestGetModelOmitsFogUUIDs(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/models/"+testModelName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		body := sampleModelJSON()
		var raw map[string]any
		if err := json.Unmarshal([]byte(body), &raw); err != nil {
			t.Fatalf("fixture JSON: %v", err)
		}
		if _, ok := raw["fogUuids"]; ok {
			t.Fatal("GET model body must not include fogUuids")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	})

	model, err := clt.GetModel(testModelName)
	if err != nil {
		t.Fatalf("GetModel: %v", err)
	}
	if model.Name != testModelName {
		t.Fatalf("name = %q", model.Name)
	}
	encoded, err := json.Marshal(model)
	if err != nil {
		t.Fatalf("marshal model: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatalf("unmarshal model: %v", err)
	}
	if _, ok := got["fogUuids"]; ok {
		t.Fatalf("encoded Model includes fogUuids: %s", encoded)
	}
}

func TestCreateModelJSON(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/models" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		req := readJSONObject(t, r)
		if req["name"] != testModelName {
			t.Fatalf("name = %v", req["name"])
		}
		if req["repo"] != "org/llama" {
			t.Fatalf("repo = %v", req["repo"])
		}
		if req["registryId"] != float64(5) {
			t.Fatalf("registryId = %v", req["registryId"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(sampleModelJSON()))
	})

	model, err := clt.CreateModel(&ModelCreateRequest{
		Name:       testModelName,
		Repo:       "org/llama",
		RegistryID: 5,
	})
	if err != nil {
		t.Fatalf("CreateModel: %v", err)
	}
	if model.UUID != testModelUUID {
		t.Fatalf("uuid = %q", model.UUID)
	}
}

func TestUpdateModel(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v3/models/"+testModelName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = readJSONObject(t, r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleModelJSON()))
	})

	repo := "org/llama-v2"
	model, err := clt.UpdateModel(testModelName, &ModelUpdateRequest{Repo: &repo})
	if err != nil {
		t.Fatalf("UpdateModel: %v", err)
	}
	if model.Name != testModelName {
		t.Fatalf("name = %q", model.Name)
	}
}

func TestDeleteModelAccepted(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v3/models/"+testModelName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	if err := clt.DeleteModel(testModelName); err != nil {
		t.Fatalf("DeleteModel: %v", err)
	}
}

func TestDeleteModelConflict(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v3/models/"+testModelName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"message":"model is bound by microservice ` + testBlockingMSUUID + `"}`))
	})

	err := clt.DeleteModel(testModelName)
	requireHTTPErrorCode(t, err, http.StatusConflict)
	if !strings.Contains(err.Error(), testBlockingMSUUID) {
		t.Fatalf("error %q does not contain microservice uuid", err)
	}
}

func TestCreateModelFromYAML(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/models/yaml" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertMultipartFormFile(t, r, "model")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleModelJSON()))
	})

	model, err := clt.CreateModelFromYAML(strings.NewReader("kind: Model\nmetadata:\n  name: llama\n"))
	if err != nil {
		t.Fatalf("CreateModelFromYAML: %v", err)
	}
	if model.Name != testModelName {
		t.Fatalf("name = %q", model.Name)
	}
}

func TestUpsertModelFromYAML(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v3/models/yaml/"+testModelName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertMultipartFormFile(t, r, "model")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleModelJSON()))
	})

	model, err := clt.UpsertModelFromYAML(testModelName, strings.NewReader("kind: Model\n"))
	if err != nil {
		t.Fatalf("UpsertModelFromYAML: %v", err)
	}
	if model.UUID != testModelUUID {
		t.Fatalf("uuid = %q", model.UUID)
	}
}

func TestGetModelLink(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/models/"+testModelName+"/link" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"fogUuids":["` + testFogUUID + `"]}`))
	})

	link, err := clt.GetModelLink(testModelName)
	if err != nil {
		t.Fatalf("GetModelLink: %v", err)
	}
	if len(link.FogUUIDs) != 1 || link.FogUUIDs[0] != testFogUUID {
		t.Fatalf("fogUuids = %v", link.FogUUIDs)
	}
}

func TestLinkModelBodyOmitsName(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/models/"+testModelName+"/link" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		req := readJSONObject(t, r)
		if _, ok := req["name"]; ok {
			t.Fatalf("link body includes name: %v", req)
		}
		uuids, ok := req["fogUuids"].([]any)
		if !ok || len(uuids) != 1 || uuids[0] != testFogUUID {
			t.Fatalf("fogUuids = %v", req["fogUuids"])
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := clt.LinkModel(testModelName, &FogLinkRequest{FogUUIDs: []string{testFogUUID}}); err != nil {
		t.Fatalf("LinkModel: %v", err)
	}
}

func TestUnlinkModelConflict(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v3/models/"+testModelName+"/link" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"message":"cannot unlink; microservice ` + testBlockingMSUUID + ` still binds llama"}`))
	})

	err := clt.UnlinkModel(testModelName, &FogLinkRequest{FogUUIDs: []string{testFogUUID}})
	requireHTTPErrorCode(t, err, http.StatusConflict)
	if !strings.Contains(err.Error(), testBlockingMSUUID) {
		t.Fatalf("error %q does not contain microservice uuid", err)
	}
}
