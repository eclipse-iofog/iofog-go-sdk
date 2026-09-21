package client

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

const testKnowledgeName = "wiki"
const testKnowledgeUUID = "b2c3d4e5-f6a7-8901-bcde-f12345678901"

func sampleKnowledgeJSON() string {
	return `{
		"uuid":"` + testKnowledgeUUID + `",
		"name":"` + testKnowledgeName + `",
		"repo":"org/wiki",
		"revision":"main",
		"registryId":5,
		"format":"dataset"
	}`
}

func TestListKnowledge(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/knowledge" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"knowledge":[` + sampleKnowledgeJSON() + `]}`))
	})

	resp, err := clt.ListKnowledge()
	if err != nil {
		t.Fatalf("ListKnowledge: %v", err)
	}
	if len(resp.Knowledge) != 1 {
		t.Fatalf("knowledge len = %d, want 1", len(resp.Knowledge))
	}
	got := resp.Knowledge[0]
	if got.Name != testKnowledgeName || got.Repo != "org/wiki" || got.RegistryID != 5 {
		t.Fatalf("unexpected knowledge: %+v", got)
	}
}

func TestGetKnowledgeOmitsFogUUIDs(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/knowledge/"+testKnowledgeName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		body := sampleKnowledgeJSON()
		var raw map[string]any
		if err := json.Unmarshal([]byte(body), &raw); err != nil {
			t.Fatalf("fixture JSON: %v", err)
		}
		if _, ok := raw["fogUuids"]; ok {
			t.Fatal("GET knowledge body must not include fogUuids")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	})

	knowledge, err := clt.GetKnowledge(testKnowledgeName)
	if err != nil {
		t.Fatalf("GetKnowledge: %v", err)
	}
	if knowledge.Name != testKnowledgeName {
		t.Fatalf("name = %q", knowledge.Name)
	}
	encoded, err := json.Marshal(knowledge)
	if err != nil {
		t.Fatalf("marshal knowledge: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatalf("unmarshal knowledge: %v", err)
	}
	if _, ok := got["fogUuids"]; ok {
		t.Fatalf("encoded Knowledge includes fogUuids: %s", encoded)
	}
}

func TestCreateKnowledgeJSON(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/knowledge" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		req := readJSONObject(t, r)
		if req["name"] != testKnowledgeName {
			t.Fatalf("name = %v", req["name"])
		}
		if req["repo"] != "org/wiki" {
			t.Fatalf("repo = %v", req["repo"])
		}
		if req["registryId"] != float64(5) {
			t.Fatalf("registryId = %v", req["registryId"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(sampleKnowledgeJSON()))
	})

	knowledge, err := clt.CreateKnowledge(&KnowledgeCreateRequest{
		Name:       testKnowledgeName,
		Repo:       "org/wiki",
		RegistryID: 5,
	})
	if err != nil {
		t.Fatalf("CreateKnowledge: %v", err)
	}
	if knowledge.UUID != testKnowledgeUUID {
		t.Fatalf("uuid = %q", knowledge.UUID)
	}
}

func TestUpdateKnowledgePointerFields(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v3/knowledge/"+testKnowledgeName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		req := readJSONObject(t, r)
		if req["repo"] != "org/wiki-v2" {
			t.Fatalf("repo = %v", req["repo"])
		}
		if _, ok := req["name"]; ok {
			t.Fatalf("PATCH body includes name: %v", req)
		}
		if _, ok := req["revision"]; ok {
			t.Fatalf("unset pointer field revision was sent: %v", req)
		}
		if _, ok := req["registryId"]; ok {
			t.Fatalf("unset pointer field registryId was sent: %v", req)
		}
		if _, ok := req["format"]; ok {
			t.Fatalf("unset pointer field format was sent: %v", req)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleKnowledgeJSON()))
	})

	repo := "org/wiki-v2"
	knowledge, err := clt.UpdateKnowledge(testKnowledgeName, &KnowledgeUpdateRequest{Repo: &repo})
	if err != nil {
		t.Fatalf("UpdateKnowledge: %v", err)
	}
	if knowledge.Name != testKnowledgeName {
		t.Fatalf("name = %q", knowledge.Name)
	}
}

func TestDeleteKnowledgeAccepted(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v3/knowledge/"+testKnowledgeName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	if err := clt.DeleteKnowledge(testKnowledgeName); err != nil {
		t.Fatalf("DeleteKnowledge: %v", err)
	}
}

func TestDeleteKnowledgeConflict(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v3/knowledge/"+testKnowledgeName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"message":"knowledge is bound by microservice ` + testBlockingMSUUID + `"}`))
	})

	err := clt.DeleteKnowledge(testKnowledgeName)
	requireHTTPErrorCode(t, err, http.StatusConflict)
	if !strings.Contains(err.Error(), testBlockingMSUUID) {
		t.Fatalf("error %q does not contain microservice uuid", err)
	}
}

func TestCreateKnowledgeFromYAML(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/knowledge/yaml" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertMultipartFormFile(t, r, "knowledge")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleKnowledgeJSON()))
	})

	knowledge, err := clt.CreateKnowledgeFromYAML(strings.NewReader("kind: Knowledge\nmetadata:\n  name: wiki\n"))
	if err != nil {
		t.Fatalf("CreateKnowledgeFromYAML: %v", err)
	}
	if knowledge.Name != testKnowledgeName {
		t.Fatalf("name = %q", knowledge.Name)
	}
}

func TestUpsertKnowledgeFromYAML(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v3/knowledge/yaml/"+testKnowledgeName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertMultipartFormFile(t, r, "knowledge")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleKnowledgeJSON()))
	})

	knowledge, err := clt.UpsertKnowledgeFromYAML(testKnowledgeName, strings.NewReader("kind: Knowledge\n"))
	if err != nil {
		t.Fatalf("UpsertKnowledgeFromYAML: %v", err)
	}
	if knowledge.UUID != testKnowledgeUUID {
		t.Fatalf("uuid = %q", knowledge.UUID)
	}
}

func TestGetKnowledgeLink(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/knowledge/"+testKnowledgeName+"/link" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"fogUuids":["` + testFogUUID + `"]}`))
	})

	link, err := clt.GetKnowledgeLink(testKnowledgeName)
	if err != nil {
		t.Fatalf("GetKnowledgeLink: %v", err)
	}
	if len(link.FogUUIDs) != 1 || link.FogUUIDs[0] != testFogUUID {
		t.Fatalf("fogUuids = %v", link.FogUUIDs)
	}
}

func TestLinkKnowledgeBodyOmitsName(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/knowledge/"+testKnowledgeName+"/link" {
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

	if err := clt.LinkKnowledge(testKnowledgeName, &FogLinkRequest{FogUUIDs: []string{testFogUUID}}); err != nil {
		t.Fatalf("LinkKnowledge: %v", err)
	}
}

func TestUnlinkKnowledgeConflict(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v3/knowledge/"+testKnowledgeName+"/link" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		req := readJSONObject(t, r)
		uuids, ok := req["fogUuids"].([]any)
		if !ok || len(uuids) != 1 || uuids[0] != testFogUUID {
			t.Fatalf("fogUuids = %v", req["fogUuids"])
		}
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"message":"cannot unlink; microservice ` + testBlockingMSUUID + ` still binds wiki"}`))
	})

	err := clt.UnlinkKnowledge(testKnowledgeName, &FogLinkRequest{FogUUIDs: []string{testFogUUID}})
	requireHTTPErrorCode(t, err, http.StatusConflict)
	if !strings.Contains(err.Error(), testBlockingMSUUID) {
		t.Fatalf("error %q does not contain microservice uuid", err)
	}
}
