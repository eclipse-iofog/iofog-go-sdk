package microservices

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadBearerToken(t *testing.T) {
	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	if err := os.WriteFile(tokenPath, []byte("  test-token  \n"), 0o600); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	token, err := readBearerToken(tokenPath)
	if err != nil {
		t.Fatalf("unexpected token read error: %v", err)
	}
	if token != "test-token" {
		t.Fatalf("unexpected token value: %q", token)
	}
}

func TestReadBearerTokenEmpty(t *testing.T) {
	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	if err := os.WriteFile(tokenPath, []byte(" \n "), 0o600); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	_, err := readBearerToken(tokenPath)
	if err == nil {
		t.Fatal("expected error for empty token file")
	}
}

func TestLoadRootCAs(t *testing.T) {
	dir := t.TempDir()
	caPath := filepath.Join(dir, "ca.crt")
	caPEM, _, _ := generateCAAndServerCert(t)
	if err := os.WriteFile(caPath, caPEM, 0o600); err != nil {
		t.Fatalf("failed to write ca file: %v", err)
	}

	pool, err := loadRootCAs(caPath)
	if err != nil {
		t.Fatalf("unexpected CA load error: %v", err)
	}
	if pool == nil {
		t.Fatal("expected non-nil cert pool")
	}
}
