package client

import (
	"net/http"
	"testing"
)

func TestPatchMicroserviceKnowledge(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v3/microservices/"+testMicroserviceUUID+"/knowledge" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		req := readJSONObject(t, r)
		if req["bindPath"] != "/knowledge" {
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
		if item["name"] != testKnowledgeName {
			t.Fatalf("item name = %v", item["name"])
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := clt.PatchMicroserviceKnowledge(testMicroserviceUUID, KnowledgeCatalog{
		BindPath:    "/knowledge",
		Permissions: "ro",
		Items:       []KnowledgeCatalogItem{{Name: testKnowledgeName}},
	})
	if err != nil {
		t.Fatalf("PatchMicroserviceKnowledge: %v", err)
	}
}
