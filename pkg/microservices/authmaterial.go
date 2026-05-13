package microservices

import (
	"crypto/x509"
	"fmt"
	"os"
	"strings"
)

func readBearerToken(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", &AuthMaterialError{Kind: "token", Path: path, Err: err}
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return "", &AuthMaterialError{Kind: "token", Path: path, Err: fmt.Errorf("token file is empty")}
	}
	return token, nil
}

func loadRootCAs(path string) (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		pool = x509.NewCertPool()
	}
	pemData, err := os.ReadFile(path)
	if err != nil {
		return nil, &AuthMaterialError{Kind: "ca", Path: path, Err: err}
	}
	if ok := pool.AppendCertsFromPEM(pemData); !ok {
		return nil, &AuthMaterialError{Kind: "ca", Path: path, Err: fmt.Errorf("failed to parse PEM certificates")}
	}
	return pool, nil
}
