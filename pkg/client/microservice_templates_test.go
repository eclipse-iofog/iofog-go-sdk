package client

import (
	"net/http"
	"strings"
	"testing"
)

const testMSTemplateName = "nginx-tpl"

func sampleMicroserviceTemplateJSON() string {
	return `{
		"name":"` + testMSTemplateName + `",
		"description":"nginx",
		"variables":[{"key":"image","defaultValue":"nginx:latest"}],
		"microservice":{"images":{"catalogItem":"nginx"}}
	}`
}

func TestListMicroserviceTemplates(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/microserviceTemplates" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"microserviceTemplates":[` + sampleMicroserviceTemplateJSON() + `]}`))
	})

	resp, err := clt.ListMicroserviceTemplates()
	if err != nil {
		t.Fatalf("ListMicroserviceTemplates: %v", err)
	}
	if len(resp.MicroserviceTemplates) != 1 {
		t.Fatalf("len = %d, want 1", len(resp.MicroserviceTemplates))
	}
	if resp.MicroserviceTemplates[0].Name != testMSTemplateName {
		t.Fatalf("name = %q", resp.MicroserviceTemplates[0].Name)
	}
}

func TestCreateMicroserviceTemplateFromYAML(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/microserviceTemplates/yaml" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertMultipartFormFile(t, r, "template")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleMicroserviceTemplateJSON()))
	})

	tpl, err := clt.CreateMicroserviceTemplateFromYAML(strings.NewReader("kind: MicroserviceTemplate\n"))
	if err != nil {
		t.Fatalf("CreateMicroserviceTemplateFromYAML: %v", err)
	}
	if tpl.Name != testMSTemplateName {
		t.Fatalf("name = %q", tpl.Name)
	}
}

func TestUpdateMicroserviceTemplateFromYAML(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v3/microserviceTemplates/yaml/"+testMSTemplateName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertMultipartFormFile(t, r, "template")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleMicroserviceTemplateJSON()))
	})

	tpl, err := clt.UpdateMicroserviceTemplateFromYAML(testMSTemplateName, strings.NewReader("kind: MicroserviceTemplate\n"))
	if err != nil {
		t.Fatalf("UpdateMicroserviceTemplateFromYAML: %v", err)
	}
	if tpl.Description != "nginx" {
		t.Fatalf("description = %q", tpl.Description)
	}
}
