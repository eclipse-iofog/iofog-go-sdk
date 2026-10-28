package client

import (
	"encoding/json"
	"testing"
)

func TestCatalogListResponseUnmarshalNumericID(t *testing.T) {
	body := []byte(`{"catalogItems":[{"id":1,"name":"Router","description":"The built-in router for Edgelet.","category":"SYSTEM","publisher":"Eclipse ioFog","picture":"none.png","isPublic":false,"registryId":1,"images":[{"containerImage":"ghcr.io/datasance/router:3.8.0-rc.1","archId":1}],"inputType":null,"outputType":null}]}`)

	var response CatalogListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("unmarshal catalog list: %v", err)
	}
	if len(response.CatalogItems) != 1 {
		t.Fatalf("expected 1 catalog item, got %d", len(response.CatalogItems))
	}
	if response.CatalogItems[0].ID != 1 {
		t.Fatalf("expected catalog id 1, got %d", response.CatalogItems[0].ID)
	}
	if response.CatalogItems[0].Name != "Router" {
		t.Fatalf("expected name Router, got %q", response.CatalogItems[0].Name)
	}
}
