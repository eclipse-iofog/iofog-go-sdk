package microservices

import "time"

// ClientOptions configures LocalAPI v3 transport/auth behavior.
type ClientOptions struct {
	Host                 string
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

// WithHost overrides LocalAPI host.
func WithHost(host string) ClientOption {
	return func(opts *ClientOptions) { opts.Host = host }
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
