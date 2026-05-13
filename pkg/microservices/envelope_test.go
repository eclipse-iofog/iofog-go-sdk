package microservices

import (
	"errors"
	"testing"
)

func TestParseV3EnvelopeSuccess(t *testing.T) {
	body := []byte(`{"success":true,"data":{"config":{"k":"v"}}}`)
	data, err := parseV3Envelope(body, 200)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	config, ok := data["config"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map config payload, got %T", data["config"])
	}
	if config["k"] != "v" {
		t.Fatalf("unexpected config value: %v", config["k"])
	}
}

func TestParseV3EnvelopeFailure(t *testing.T) {
	body := []byte(`{"success":false,"error":{"code":"UNAUTHORIZED","message":"invalid JWT token","details":{"a":"b"}}}`)
	_, err := parseV3Envelope(body, 401)
	if err == nil {
		t.Fatal("expected parse error for failed envelope")
	}
	var apiErr *V3APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected V3APIError, got %T", err)
	}
	if apiErr.Code != "UNAUTHORIZED" {
		t.Fatalf("unexpected code: %s", apiErr.Code)
	}
	if apiErr.StatusCode != 401 {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
}
