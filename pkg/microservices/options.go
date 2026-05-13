package microservices

import (
	"strings"
	"time"
)

// ClientOptions configures LocalAPI v3 transport/auth behavior.
type ClientOptions struct {
	Host                 string
	FallbackHosts        []string
	Port                 int
	UseTLS               bool
	TokenPath            string
	CAPath               string
	RequestTimeout       time.Duration
	WSHandshakeTimeout   time.Duration
	WSReconnectBaseDelay time.Duration
	WSReconnectMaxDelay  time.Duration
}

// ClientOption mutates ClientOptions.
type ClientOption func(*ClientOptions)

func defaultClientOptions() ClientOptions {
	return ClientOptions{
		Host:                 HostDefault,
		FallbackHosts:        []string{FallbackHostLocal},
		Port:                 PortIoFog,
		UseTLS:               SSLDefault,
		TokenPath:            DefaultServiceAccountTokenPath,
		CAPath:               DefaultServiceAccountCAPath,
		RequestTimeout:       DefaultRequestTimeout,
		WSHandshakeTimeout:   DefaultWSHandshakeTimeout,
		WSReconnectBaseDelay: DefaultWSReconnectBaseDelay,
		WSReconnectMaxDelay:  DefaultWSReconnectMaxDelay,
	}
}

func applyClientOptions(base ClientOptions, opts ...ClientOption) ClientOptions {
	out := base
	for _, opt := range opts {
		if opt != nil {
			opt(&out)
		}
	}
	if out.Host == "" {
		out.Host = HostDefault
	}
	out.FallbackHosts = sanitizeFallbackHosts(out.Host, out.FallbackHosts)
	if out.Port <= 0 {
		out.Port = PortIoFog
	}
	if out.TokenPath == "" {
		out.TokenPath = DefaultServiceAccountTokenPath
	}
	if out.CAPath == "" {
		out.CAPath = DefaultServiceAccountCAPath
	}
	if out.RequestTimeout <= 0 {
		out.RequestTimeout = DefaultRequestTimeout
	}
	if out.WSHandshakeTimeout <= 0 {
		out.WSHandshakeTimeout = DefaultWSHandshakeTimeout
	}
	if out.WSReconnectBaseDelay <= 0 {
		out.WSReconnectBaseDelay = DefaultWSReconnectBaseDelay
	}
	if out.WSReconnectMaxDelay <= 0 {
		out.WSReconnectMaxDelay = DefaultWSReconnectMaxDelay
	}
	if out.WSReconnectMaxDelay < out.WSReconnectBaseDelay {
		out.WSReconnectMaxDelay = out.WSReconnectBaseDelay
	}
	return out
}

func sanitizeFallbackHosts(primary string, hosts []string) []string {
	seen := make(map[string]struct{}, len(hosts)+1)
	normalizedPrimary := strings.TrimSpace(primary)
	if normalizedPrimary != "" {
		seen[normalizedPrimary] = struct{}{}
	}

	out := make([]string, 0, len(hosts))
	for _, host := range hosts {
		normalized := strings.TrimSpace(host)
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	return out
}

// WithHost overrides LocalAPI host.
func WithHost(host string) ClientOption {
	return func(opts *ClientOptions) { opts.Host = host }
}

// WithFallbackHosts overrides fallback LocalAPI hosts.
func WithFallbackHosts(hosts ...string) ClientOption {
	return func(opts *ClientOptions) { opts.FallbackHosts = hosts }
}

// WithPort overrides LocalAPI port.
func WithPort(port int) ClientOption {
	return func(opts *ClientOptions) { opts.Port = port }
}

// WithTLS toggles HTTPS/WSS.
func WithTLS(enabled bool) ClientOption {
	return func(opts *ClientOptions) { opts.UseTLS = enabled }
}

// WithTokenPath overrides service-account token path.
func WithTokenPath(path string) ClientOption {
	return func(opts *ClientOptions) { opts.TokenPath = path }
}

// WithCAPath overrides service-account CA path.
func WithCAPath(path string) ClientOption {
	return func(opts *ClientOptions) { opts.CAPath = path }
}

// WithRequestTimeout sets HTTP request timeout.
func WithRequestTimeout(timeout time.Duration) ClientOption {
	return func(opts *ClientOptions) { opts.RequestTimeout = timeout }
}

// WithWSHandshakeTimeout sets websocket handshake timeout.
func WithWSHandshakeTimeout(timeout time.Duration) ClientOption {
	return func(opts *ClientOptions) { opts.WSHandshakeTimeout = timeout }
}

// WithWSReconnectDelays sets websocket reconnect backoff bounds.
func WithWSReconnectDelays(base, max time.Duration) ClientOption {
	return func(opts *ClientOptions) {
		opts.WSReconnectBaseDelay = base
		opts.WSReconnectMaxDelay = max
	}
}
