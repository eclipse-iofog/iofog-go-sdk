package microservices

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (opts ClientOptions) restBaseURL() string {
	return opts.restBaseURLForHost(opts.Host)
}

func (opts ClientOptions) restBaseURLForHost(host string) string {
	scheme := SchemeHTTP
	if opts.UseTLS {
		scheme = SchemeHTTPS
	}
	return fmt.Sprintf("%s://%s:%d", scheme, host, opts.Port)
}

func (opts ClientOptions) wsBaseURL() string {
	return opts.wsBaseURLForHost(opts.Host)
}

func (opts ClientOptions) wsBaseURLForHost(host string) string {
	scheme := SchemeWS
	if opts.UseTLS {
		scheme = SchemeWSS
	}
	return fmt.Sprintf("%s://%s:%d", scheme, host, opts.Port)
}

func (opts ClientOptions) requestHosts() []string {
	out := make([]string, 0, 1+len(opts.FallbackHosts))
	out = append(out, opts.Host)
	out = append(out, opts.FallbackHosts...)
	return out
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

func isRetriableHostError(err error) bool {
	if err == nil {
		return false
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return isRetriableHostError(urlErr.Err)
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "dial tcp") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "i/o timeout") ||
		strings.Contains(msg, "connection refused")
}
