package k8s

import (
	"testing"
)

func TestCreation(t *testing.T) {
	// Here just to test compilation
	client := &Client{}
	if client == nil {
		t.Error("This is impossible")
	}
}
