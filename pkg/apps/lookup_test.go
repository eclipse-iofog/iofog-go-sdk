package apps

import (
	"errors"
	"testing"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
)

func TestLookupErrorAllowNotFound(t *testing.T) {
	t.Parallel()

	if err := lookupErrorAllowNotFound(nil); err != nil {
		t.Fatalf("expected nil for nil input, got %v", err)
	}

	notFound := client.NewNotFoundError("missing")
	if err := lookupErrorAllowNotFound(notFound); err != nil {
		t.Fatalf("expected nil for NotFound, got %v", err)
	}

	other := client.NewHTTPError("server error", 500)
	if err := lookupErrorAllowNotFound(other); err == nil {
		t.Fatal("expected error for non-NotFound lookup failure")
	} else if !errors.Is(err, other) {
		t.Fatalf("expected original error, got %v", err)
	}
}
