package microservices

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func (opts ClientOptions) restBaseURL() string {
	scheme := SchemeHTTP
	if opts.UseTLS {
		scheme = SchemeHTTPS
	}
	return fmt.Sprintf("%s://%s:%d", scheme, opts.Host, opts.Port)
}

func (opts ClientOptions) wsBaseURL() string {
	scheme := SchemeWS
	if opts.UseTLS {
		scheme = SchemeWSS
	}
	return fmt.Sprintf("%s://%s:%d", scheme, opts.Host, opts.Port)
}

func buildHTTPClient(opts ClientOptions) (*http.Client, error) {
	transport, err := buildHTTPTransport(opts)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Timeout:   opts.RequestTimeout,
		Transport: transport,
	}, nil
}

func buildHTTPTransport(opts ClientOptions) (*http.Transport, error) {
	transport := &http.Transport{}
	if !opts.UseTLS {
		return transport, nil
	}
	rootCAs, err := loadRootCAs(opts.CAPath)
	if err != nil {
		return nil, err
	}
	transport.TLSClientConfig = &tls.Config{
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS12,
	}
	return transport, nil
}

func sleepWithCap(attempt int, baseDelay, maxDelay time.Duration) {
	if attempt < 0 {
		attempt = 0
	}
	delay := baseDelay
	for i := 0; i < attempt; i++ {
		delay *= 2
		if delay >= maxDelay {
			delay = maxDelay
			break
		}
	}
	time.Sleep(delay)
}
