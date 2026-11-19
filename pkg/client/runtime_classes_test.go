package client

import (
	"net/http"
	"strings"
	"testing"
)

const testRuntimeClassName = "nvidia"
const testRuntimeHandler = "nvidia"

func sampleRuntimeClassJSON() string {
	return `{"name":"` + testRuntimeClassName + `","handler":"` + testRuntimeHandler + `"}`
}

func TestListRuntimeClasses(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/runtimeClasses" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"runtimeClasses":[` + sampleRuntimeClassJSON() + `]}`))
	})

	resp, err := clt.ListRuntimeClasses()
	if err != nil {
		t.Fatalf("ListRuntimeClasses: %v", err)
	}
	if len(resp.RuntimeClasses) != 1 || resp.RuntimeClasses[0].Name != testRuntimeClassName {
		t.Fatalf("unexpected list: %+v", resp.RuntimeClasses)
	}
}

func TestGetRuntimeClass(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3/runtimeClasses/"+testRuntimeClassName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleRuntimeClassJSON()))
	})

	rc, err := clt.GetRuntimeClass(testRuntimeClassName)
	if err != nil {
		t.Fatalf("GetRuntimeClass: %v", err)
	}
	if rc.Handler != testRuntimeHandler {
		t.Fatalf("handler = %q", rc.Handler)
	}
}

func TestCreateRuntimeClass(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/runtimeClasses" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		req := readJSONObject(t, r)
		if req["name"] != testRuntimeClassName || req["handler"] != testRuntimeHandler {
			t.Fatalf("body = %v", req)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(sampleRuntimeClassJSON()))
	})

	rc, err := clt.CreateRuntimeClass(&RuntimeClassCreateRequest{
		Name:    testRuntimeClassName,
		Handler: testRuntimeHandler,
	})
	if err != nil {
		t.Fatalf("CreateRuntimeClass: %v", err)
	}
	if rc.Name != testRuntimeClassName {
		t.Fatalf("name = %q", rc.Name)
	}
}

func TestUpdateRuntimeClass(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v3/runtimeClasses/"+testRuntimeClassName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		_ = readJSONObject(t, r)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleRuntimeClassJSON()))
	})

	rc, err := clt.UpdateRuntimeClass(testRuntimeClassName, &RuntimeClassUpdateRequest{Handler: "nvidia-cdi"})
	if err != nil {
		t.Fatalf("UpdateRuntimeClass: %v", err)
	}
	if rc.Name != testRuntimeClassName {
		t.Fatalf("name = %q", rc.Name)
	}
}

func TestDeleteRuntimeClass(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v3/runtimeClasses/"+testRuntimeClassName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	if err := clt.DeleteRuntimeClass(testRuntimeClassName); err != nil {
		t.Fatalf("DeleteRuntimeClass: %v", err)
	}
}

func TestCreateRuntimeClassFromYAML(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/runtimeClasses/yaml" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertMultipartFormFile(t, r, "runtimeClass")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleRuntimeClassJSON()))
	})

	rc, err := clt.CreateRuntimeClassFromYAML(strings.NewReader("kind: RuntimeClass\n"))
	if err != nil {
		t.Fatalf("CreateRuntimeClassFromYAML: %v", err)
	}
	if rc.Name != testRuntimeClassName {
		t.Fatalf("name = %q", rc.Name)
	}
}

func TestUpsertRuntimeClassFromYAML(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v3/runtimeClasses/yaml/"+testRuntimeClassName {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		assertMultipartFormFile(t, r, "runtimeClass")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sampleRuntimeClassJSON()))
	})

	rc, err := clt.UpsertRuntimeClassFromYAML(testRuntimeClassName, strings.NewReader("kind: RuntimeClass\n"))
	if err != nil {
		t.Fatalf("UpsertRuntimeClassFromYAML: %v", err)
	}
	if rc.Handler != testRuntimeHandler {
		t.Fatalf("handler = %q", rc.Handler)
	}
}

func TestLinkRuntimeClass(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/runtimeClasses/"+testRuntimeClassName+"/link" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		req := readJSONObject(t, r)
		if _, ok := req["name"]; ok {
			t.Fatalf("link body includes name: %v", req)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	if err := clt.LinkRuntimeClass(testRuntimeClassName, &FogLinkRequest{FogUUIDs: []string{testFogUUID}}); err != nil {
		t.Fatalf("LinkRuntimeClass: %v", err)
	}
}

func TestLinkRuntimeClassNonEdgeletFog(t *testing.T) {
	clt := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v3/runtimeClasses/"+testRuntimeClassName+"/link" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"fog ` + testFogUUID + ` is not containerEngine edgelet"}`))
	})

	err := clt.LinkRuntimeClass(testRuntimeClassName, &FogLinkRequest{FogUUIDs: []string{testFogUUID}})
	requireHTTPErrorCode(t, err, http.StatusBadRequest)
}
