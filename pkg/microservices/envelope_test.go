package microservices

import (
	"errors"
	"testing"
)

func TestParseEdgeletAPIEnvelopeSuccess(t *testing.T) {
	body := []byte(`{"success":true,"data":{"config":{"k":"v"}}}`)
	data, err := parseEdgeletAPIEnvelope(body, 200)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	config, ok := data["config"].(map[string]any)
	if !ok {
		t.Fatalf("expected map config payload, got %T", data["config"])
	}
	if config["k"] != "v" {
		t.Fatalf("unexpected config value: %v", config["k"])
	}
}

func TestParseEdgeletAPIEnvelopeFailure(t *testing.T) {
	body := []byte(`{"success":false,"error":{"code":"UNAUTHORIZED","message":"invalid JWT token","details":{"a":"b"}}}`)
	_, err := parseEdgeletAPIEnvelope(body, 401)
	if err == nil {
		t.Fatal("expected parse error for failed envelope")
	}
	var apiErr *EdgeletAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected EdgeletAPIError, got %T", err)
	}
	if apiErr.Code != "UNAUTHORIZED" {
		t.Fatalf("unexpected code: %s", apiErr.Code)
	}
	if apiErr.StatusCode != 401 {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
}
