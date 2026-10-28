package microservices

import (
	"crypto/x509"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func readFileAtPath(path string) ([]byte, error) {
	clean := filepath.Clean(path)
	dir, base := filepath.Split(clean)
	if base == "" {
		return nil, errors.New("invalid file path")
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	f, err := root.Open(base)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func readBearerToken(path string) (string, error) {
	data, err := readFileAtPath(path)
	if err != nil {
		return "", &AuthMaterialError{Kind: "token", Path: path, Err: err}
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", &AuthMaterialError{Kind: "token", Path: path, Err: errors.New("token file is empty")}
	}
	return token, nil
}

func loadRootCAs(path string) (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		pool = x509.NewCertPool()
	}
	pemData, err := readFileAtPath(path)
	if err != nil {
		return nil, &AuthMaterialError{Kind: "ca", Path: path, Err: err}
	}
	if ok := pool.AppendCertsFromPEM(pemData); !ok {
		return nil, &AuthMaterialError{Kind: "ca", Path: path, Err: errors.New("failed to parse PEM certificates")}
	}
	return pool, nil
}
